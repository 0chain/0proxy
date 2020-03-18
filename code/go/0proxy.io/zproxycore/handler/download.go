package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"

	"0proxy.io/zproxycore/common"
	"github.com/0chain/gosdk/zboxcore/sdk"
)

func download(w http.ResponseWriter, r *http.Request) {
	allocation := r.FormValue("allocation")
	clientJSON := r.FormValue("client_json")
	remotePath := r.FormValue("remote_path")
	authTicket := r.FormValue("auth_ticket")
	numBlocks := r.FormValue("blocks_per_marker")
	numBlocksInt, _ := strconv.Atoi(numBlocks)
	if numBlocksInt == 0 {
		numBlocksInt = 10
	}
	if len(remotePath) == 0 && len(authTicket) == 0 {
		handleError(w, fmt.Errorf("invalid_params", "remote_path / auth_ticket missing"))
	}

	err := initSDK(clientJSON)
	if err != nil {
		handleError(w, fmt.Errorf("skd_not_initialized", "Unable to initialize sdk"))
	}
	sdk.SetNumBlockDownloads(numBlocksInt)

	var at *sdk.AuthTicket
	downloadUsingAT := false
	if len(authTicket) > 0 {
		downloadUsingAT = true
		at = sdk.InitAuthTicket(authTicket)
	}

	var localFilePath string
	wg := &sync.WaitGroup{}
	statusBar := &StatusBar{wg: wg}
	wg.Add(1)
	if downloadUsingAT {
		allocationObj, err := sdk.GetAllocationFromAuthTicket(authTicket)
		if err != nil {
			handleError(w, err)
		}

		fileName, err := at.GetFileName()
		if err != nil {
			handleError(w, err)
		}
		localFilePath = common.GetPath(allocationObj.ID, fileName)
		lookuphash, err := at.GetLookupHash()
		if err != nil {
			handleError(w, err)
		}

		err = allocationObj.DownloadFromAuthTicket(localFilePath, authTicket, lookuphash, fileName, statusBar)
		if err != nil {
			handleError(w, err)
		}
	} else {
		fileName := filepath.Base(remotePath)
		localFilePath = common.GetPath(allocation, fileName)

		allocationObj, err := sdk.GetAllocation(allocation)
		if err != nil {
			handleError(w, err)
		}

		err = allocationObj.DownloadFile(localFilePath, remotePath, statusBar)
		if err != nil {
			handleError(w, err)
		}
	}
	if !statusBar.success {
		wg.Wait()
		handleError(w, fmt.Errorf("Upload not successfull"))
	}
	http.ServeFile(w, r, localFilePath)
	common.DeleletFile(localFilePath)
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
