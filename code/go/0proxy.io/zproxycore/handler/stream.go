package handler

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"0proxy.io/core/common"
	. "0proxy.io/core/logging"
	"github.com/0chain/gosdk/zboxcore/sdk"
	"go.uber.org/zap"
)

// Download is to download a file from dStorage
func Stream(w http.ResponseWriter, r *http.Request) (string, error) {
	if r.Method != http.MethodGet {
		return "", methodError("stream", http.MethodGet)
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

	inputRanges := r.Header.Get("Range")

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

	var allocationObj *sdk.Allocation
	var metaData *sdk.ConsolidatedFileMeta
	var lookuphash, localFilePath, fileName string
	if downloadUsingAT {
		allocationObj, err = sdk.GetAllocationFromAuthTicket(authTicket)
		if err != nil {
			return "", common.NewError("get_allocation_failed", err.Error())
		}

		lookuphash, err = at.GetLookupHash()
		if err != nil {
			return "", common.NewError("get_lookuphash_failed", err.Error())
		}

		metaData, err = allocationObj.GetFileMetaFromAuthTicket(authTicket, lookuphash)
		if err != nil {
			return "", common.NewError("get_filemeta_failed", err.Error())
		}

		fileName, err = at.GetFileName()
		if err != nil {
			return "", common.NewError("get_file_name_failed", err.Error())
		}
	} else {
		allocationObj, err = sdk.GetAllocation(allocation)
		if err != nil {
			return "", common.NewError("get_allocation_failed", err.Error())
		}

		metaData, err = allocationObj.GetFileMeta(remotePath)
		if err != nil {
			return "", common.NewError("get_filemeta_failed", err.Error())
		}

		fileName = filepath.Base(remotePath)
	}

	startBlockInt, endBlockInt, chunkStart, chunkEnd := calChunk(allocationObj, numBlocksInt, int(metaData.Size), inputRanges)

	createDirIfNotExists(allocationObj.ID)
	localFilePath = getPathForStream(allocationObj.ID, fileName, startBlockInt, endBlockInt)
	if _, err := os.Stat(localFilePath); os.IsNotExist(err) {
		wg := &sync.WaitGroup{}
		statusBar := &StatusBar{wg: wg}
		wg.Add(1)
		if downloadUsingAT {
			rxPay, _ := strconv.ParseBool(r.FormValue("rx_pay"))
			Logger.Info("Doing file download using authTicket", zap.Any("filename", fileName), zap.Any("allocation", allocationObj.ID), zap.Any("lookuphash", lookuphash), zap.Any("startblock", startBlockInt), zap.Any("endblock", endBlockInt))
			err = allocationObj.DownloadFromAuthTicketByBlocks(localFilePath, authTicket, int64(startBlockInt), int64(endBlockInt), numBlocksInt, lookuphash, fileName, rxPay, statusBar)
			if err != nil {
				return "", common.NewError("download_from_auth_ticket_failed", err.Error())
			}
		} else {
			Logger.Info("Doing file download", zap.Any("remotepath", remotePath), zap.Any("allocation", allocation), zap.Any("startblock", startBlockInt), zap.Any("endblock", endBlockInt))
			err = allocationObj.DownloadFileByBlock(localFilePath, remotePath, int64(startBlockInt), int64(endBlockInt), numBlocksInt, statusBar)
			if err != nil {
				return "", common.NewError("download_file_failed", err.Error())
			}
		}
		wg.Wait()
		if !statusBar.success {
			return "", common.NewError("download_status_failed", "Status bar not success")
		}
	}

	w.Header().Set("Content-Length", strconv.Itoa(chunkEnd-chunkStart+1))
	contentRange := fmt.Sprintf("bytes %d-%d/%d", chunkStart, chunkEnd, metaData.Size)
	w.Header().Set("Content-Type", metaData.MimeType)
	w.Header().Set("Content-Range", contentRange)
	return localFilePath, nil
}

func calChunk(alloc *sdk.Allocation, numBlocks int, fileSize int, inputRange string) (int, int, int, int) {
	dataShards := alloc.DataShards
	chunkMultiplier := 1
	if dataShards > 0 {
		for dataShards*chunkMultiplier < numBlocks {
			chunkMultiplier++
		}
	}

	chunkSize := chunkMultiplier * dataShards * 65536
	ranges := calRange(inputRange, fileSize)

	chunkStart := ranges[0]
	chunkEnd := (int(math.Floor(float64((chunkStart+chunkSize)/chunkSize))) * chunkSize) - 1
	if chunkEnd >= fileSize {
		chunkEnd = fileSize - 1
	}

	if ranges[0] == 0 && ranges[1] == 1 {
		chunkEnd = 1
	}

	chunkNo := int(math.Floor(float64(chunkStart/chunkSize)) + 1)
	return ((chunkNo - 1) * chunkMultiplier) + 1, ((chunkNo - 1) * chunkMultiplier) + chunkMultiplier, chunkStart, chunkEnd
}

func calRange(inputRange string, fileSize int) []int {
	ranges := make([]int, 2)
	httpRanges, err := parseRange(inputRange, int64(fileSize))
	if err != nil || len(httpRanges) == 0 {
		ranges[0] = 0
		ranges[1] = fileSize - 1
		return ranges
	}

	ranges[0] = int(httpRanges[0].start)
	if len(httpRanges) < 2 {
		ranges[1] = fileSize - 1
	} else {
		ranges[1] = int(httpRanges[1].start)
	}
	return ranges
}
