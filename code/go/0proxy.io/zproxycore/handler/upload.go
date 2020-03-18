package handler

import (
	"context"
	"net/http"
	"strconv"
	"sync"

	"0proxy.io/zproxycore/common"
	"github.com/0chain/gosdk/zboxcore/sdk"
)

func Upload(ctx context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodPost || r.Method != http.MethodPut {
		return nil, common.NewError("invalid_method", "Invalid method for upload endpoint, Use POST or PUT.")
	}
	allocation := r.FormValue("allocation")
	clientJSON := r.FormValue("client_json")
	remotePath := r.FormValue("remote_path")
	encrypt := r.FormValue("encrypt")
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	common.CreateDirIfNotExists(allocation)

	localFilePath, err := common.WriteFile(file, common.GetPath(allocation, fileHeader.Filename))
	if err != nil {
		return nil, err
	}

	err = initSDK(clientJSON)
	if err != nil {
		return nil, common.NewError("sdk_not_initialized", "Unable to initialize gosdk with the given client details")
	}

	allocationObj, err := sdk.GetAllocation(allocation)
	if err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}
	statusBar := &StatusBar{wg: wg}
	wg.Add(1)
	if r.Method == http.MethodPost {
		encryptBool, _ := strconv.ParseBool(encrypt)
		if encryptBool {
			err = allocationObj.EncryptAndUploadFile(localFilePath, remotePath, statusBar)
		} else {
			err = allocationObj.UploadFile(localFilePath, remotePath, statusBar)
		}
	} else {
		err = allocationObj.UpdateFile(localFilePath, remotePath, statusBar)
	}
	if err != nil {
		return nil, common.NewError("upload_failed", "Upload failed")
	}

	wg.Wait()
	if !statusBar.success {
		return nil, common.NewError("upload_failed", "Upload failed")
	}
	err = common.DeleletFile(localFilePath)
	var response struct {
		Message string `json:"msg"`
	}
	response.Message = "Upload done"
	return response, err
}
