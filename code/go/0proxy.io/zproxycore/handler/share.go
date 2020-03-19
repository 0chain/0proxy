package handler

import (
	"context"
	"net/http"
	"path/filepath"

	"0proxy.io/core/common"
	"github.com/0chain/gosdk/zboxcore/fileref"
	"github.com/0chain/gosdk/zboxcore/sdk"
)

// Share is to share file in dStorage
func Share(ctx context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodGet {
		return nil, methodError("share", http.MethodGet)
	}

	allocation := r.FormValue("allocation")
	clientJSON := r.FormValue("client_json")
	err := validateClientDetails(allocation, clientJSON)
	if err != nil {
		return nil, err
	}

	remotePath := r.FormValue("remote_path")
	if len(remotePath) == 0 {
		return nil, common.NewError("invalid_param", "Please provide remote_path for share")
	}

	refereeClientID := r.FormValue("referee_client_id")
	encryptionpublickey := r.FormValue(("encryption_public_key"))

	err = initSDK(clientJSON)
	if err != nil {
		return nil, common.NewError("sdk_not_initialized", "Unable to initialize gosdk with the given client details")
	}

	allocationObj, err := sdk.GetAllocation(allocation)
	if err != nil {
		return nil, common.NewError("get_allocation_failed", err.Error())

	}

	refType := fileref.FILE
	statsMap, err := allocationObj.GetFileStats(remotePath)
	if err != nil {
		return nil, common.NewError("get_file_stats_failed", err.Error())
	}

	isFile := false
	for _, v := range statsMap {
		if v != nil {
			isFile = true
			break
		}
	}
	if !isFile {
		refType = fileref.DIRECTORY
	}

	var fileName string
	_, fileName = filepath.Split(remotePath)

	at, err := allocationObj.GetAuthTicket(remotePath, fileName, refType, refereeClientID, encryptionpublickey)
	if err != nil {
		return nil, common.NewError("get_auth_ticket_failed", err.Error())
	}

	var response struct {
		AuthTicket string `json:"auth_ticket"`
	}
	response.AuthTicket = at
	return response, nil
}
