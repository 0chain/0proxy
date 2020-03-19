package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"

	"0proxy.io/core/common"
	. "0proxy.io/core/logging"
	"github.com/0chain/gosdk/zboxcore/sdk"
	"go.uber.org/zap"
)

// Download is to download a file from dStorage
func Download(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, methodError("download", http.MethodGet))
	}

	allocation := r.FormValue("allocation")
	clientJSON := r.FormValue("client_json")
	err := validateClientDetails(allocation, clientJSON)
	if err != nil {
		handleError(w, err)
	}

	remotePath := r.FormValue("remote_path")
	authTicket := r.FormValue("auth_ticket")
	if len(remotePath) == 0 && len(authTicket) == 0 {
		handleError(w, common.NewError("invalid_params", "Please provide remote_path OR auth_ticket to download"))
	}

	numBlocks := r.FormValue("blocks_per_marker")
	numBlocksInt, _ := strconv.Atoi(numBlocks)
	if numBlocksInt == 0 {
		numBlocksInt = 10
	}

	err = initSDK(clientJSON)
	if err != nil {
		handleError(w, common.NewError("sdk_not_initialized", "Unable to initialize gosdk with the given client details"))
	}
	sdk.SetNumBlockDownloads(numBlocksInt)

	var at *sdk.AuthTicket
	downloadUsingAT := false
	if len(authTicket) > 0 {
		downloadUsingAT = true
		at = sdk.InitAuthTicket(authTicket)
	}

	var localFilePath, fileName string
	wg := &sync.WaitGroup{}
	statusBar := &StatusBar{wg: wg}
	wg.Add(1)
	if downloadUsingAT {
		allocationObj, err := sdk.GetAllocationFromAuthTicket(authTicket)
		if err != nil {
			handleError(w, common.NewError("get_allocation_failed", err.Error()))
		}

		fileName, err = at.GetFileName()
		if err != nil {
			handleError(w, common.NewError("get_file_name_failed", err.Error()))
		}

		createDirIfNotExists(allocationObj.ID)
		localFilePath = getPath(allocationObj.ID, fileName)
		lookuphash, err := at.GetLookupHash()
		if err != nil {
			handleError(w, common.NewError("get_lookuphash_failed", err.Error()))
		}

		Logger.Info("Doing file download using authTicket", zap.Any("filename", fileName), zap.Any("allocation", allocationObj.ID), zap.Any("lookuphash", lookuphash))
		err = allocationObj.DownloadFromAuthTicket(localFilePath, authTicket, lookuphash, fileName, statusBar)
		if err != nil {
			handleError(w, common.NewError("download_from_auth_ticket_failed", err.Error()))
		}
	} else {
		createDirIfNotExists(allocation)
		fileName = filepath.Base(remotePath)
		localFilePath = getPath(allocation, fileName)

		allocationObj, err := sdk.GetAllocation(allocation)
		if err != nil {
			handleError(w, common.NewError("get_allocation_failed", err.Error()))
		}

		Logger.Info("Doing file download", zap.Any("remotepath", remotePath), zap.Any("allocation", allocationObj.ID))
		err = allocationObj.DownloadFile(localFilePath, remotePath, statusBar)
		if err != nil {
			handleError(w, common.NewError("download_file_failed", err.Error()))
		}
	}
	wg.Wait()
	if !statusBar.success {
		handleError(w, common.NewError("download_status_failed", "Status bar not success"))
	}
	w.Header().Add("Content-Disposition", fmt.Sprintf("Attachment; filename=%s", fileName))
	http.ServeFile(w, r, localFilePath)
	deleletFile(localFilePath)
}

func handleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	data := make(map[string]interface{}, 2)
	data["error"] = err.Error()
	if cerr, ok := err.(*common.Error); ok {
		data["code"] = cerr.Code
	}
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(data)
	http.Error(w, buf.String(), 400)
}
