package handler

import (
	"context"
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
func Download(ctx context.Context, r *http.Request) (string, error) {
	if r.Method != http.MethodGet {
		return "", methodError("download", http.MethodGet)
	}

	allocation := r.FormValue("allocation")
	clientJSON := r.FormValue("client_json")
	err := validateClientDetails(allocation, clientJSON)
	if err != nil {
		return "", err
	}

	remotePath := r.FormValue("remote_path")
	authTicket := r.FormValue("auth_ticket")
	if len(remotePath) == 0 && len(authTicket) == 0 {
		return "", common.NewError("invalid_params", "Please provide remote_path OR auth_ticket to download")
	}

	numBlocks := r.FormValue("blocks_per_marker")
	numBlocksInt, _ := strconv.Atoi(numBlocks)
	if numBlocksInt == 0 {
		numBlocksInt = 10
	}

	err = initSDK(clientJSON)
	if err != nil {
		return "", common.NewError("sdk_not_initialized", "Unable to initialize gosdk with the given client details")
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
		rxPay, _ := strconv.ParseBool(r.FormValue("rx_pay"))
		allocationObj, err := sdk.GetAllocationFromAuthTicket(authTicket)
		if err != nil {
			return "", common.NewError("get_allocation_failed", err.Error())
		}
		fileName := r.FormValue("file_name")
		if len(fileName) == 0 {
			fileName, err = at.GetFileName()
			if err != nil {
				return "", common.NewError("get_file_name_failed", err.Error())
			}
		}

		createDirIfNotExists(allocationObj.ID)
		localFilePath = getPath(allocationObj.ID, fileName)
		deleletFile(localFilePath)
		lookuphash := r.FormValue("lookup_hash")
		if len(lookuphash) == 0 {
			lookuphash, err = at.GetLookupHash()
			if err != nil {
				return "", common.NewError("get_lookuphash_failed", err.Error())
			}
		}

		Logger.Info("Doing file download using authTicket", zap.Any("filename", fileName), zap.Any("allocation", allocationObj.ID), zap.Any("lookuphash", lookuphash))
		err = allocationObj.DownloadFromAuthTicket(localFilePath, authTicket, lookuphash, fileName, rxPay, statusBar)
		if err != nil {
			return "", common.NewError("download_from_auth_ticket_failed", err.Error())
		}
	} else {
		createDirIfNotExists(allocation)
		fileName = filepath.Base(remotePath)
		localFilePath = getPath(allocation, fileName)
		deleletFile(localFilePath)

		allocationObj, err := sdk.GetAllocation(allocation)
		if err != nil {
			return "", common.NewError("get_allocation_failed", err.Error())
		}

		Logger.Info("Doing file download", zap.Any("remotepath", remotePath), zap.Any("allocation", allocation))
		err = allocationObj.DownloadFile(localFilePath, remotePath, statusBar)
		if err != nil {
			return "", common.NewError("download_file_failed", err.Error())
		}
	}
	wg.Wait()
	if !statusBar.success {
		return "", statusBar.err
	}

	return localFilePath, nil
}
