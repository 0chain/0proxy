package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"0proxy.io/core/common"
	. "0proxy.io/core/logging"
	"github.com/0chain/gosdk/zboxcore/fileref"
	"github.com/0chain/gosdk/zboxcore/sdk"
	"go.uber.org/zap"
)

// Upload is to upload file to dStorage
func Upload(ctx context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		return nil, methodError("upload", http.MethodPost)
	}

	allocation := r.FormValue("allocation")
	clientJSON := r.FormValue("client_json")
	err := validateClientDetails(allocation, clientJSON)
	if err != nil {
		return nil, err
	}

	remotePath := r.FormValue("remote_path")
	if len(remotePath) == 0 {
		return nil, common.NewError("invalid_param", "Please provide remote_path for upload")
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, common.NewError("invalid_params", "Unable to get file for upload :"+err.Error())
	}
	defer file.Close()
	encrypt := r.FormValue("encrypt")

	createDirIfNotExists(allocation)

	localFilePath, err := writeFile(file, getPath(allocation, fileHeader.Filename))
	if err != nil {
		return nil, common.NewError("write_local_temp_file_failed", err.Error())
	}
	defer deleletFile(localFilePath)

	err = initSDK(clientJSON)
	if err != nil {
		return nil, common.NewError("sdk_not_initialized", "Unable to initialize gosdk with the given client details")
	}

	allocationObj, err := sdk.GetAllocation(allocation)
	if err != nil {
		return nil, common.NewError("get_allocation_failed", err.Error())
	}

	fileAttrs := r.FormValue("file_attrs")
	var attrs fileref.Attributes
	if len(fileAttrs) > 0 {
		err := json.Unmarshal([]byte(fileAttrs), &attrs)
		if err != nil {
			return nil, common.NewError("failed_to_parse_file_attrs", err.Error())
		}
	}

	wg := &sync.WaitGroup{}
	statusBar := &StatusBar{wg: wg}
	wg.Add(1)
	if r.Method == http.MethodPost {
		encryptBool, _ := strconv.ParseBool(encrypt)
		if encryptBool {
			Logger.Info("Doing encrypted file upload with", zap.Any("remotepath", remotePath), zap.Any("allocation", allocationObj.ID))
			err = allocationObj.EncryptAndUploadFile(localFilePath, remotePath, attrs, statusBar)
		} else {
			Logger.Info("Doing file upload with", zap.Any("remotepath", remotePath), zap.Any("allocation", allocationObj.ID))
			err = allocationObj.UploadFile(localFilePath, remotePath, attrs, statusBar)
		}
	} else {
		Logger.Info("Doing file update with", zap.Any("remotepath", remotePath), zap.Any("allocation", allocationObj.ID))
		err = allocationObj.UpdateFile(localFilePath, remotePath, attrs, statusBar)
	}
	if err != nil {
		return nil, common.NewError("upload_file_failed", err.Error())
	}

	wg.Wait()
	if !statusBar.success {
		return nil, statusBar.err
	}

	resp := response{
		Message: "Upload done successfully",
	}
	return resp, nil
}
