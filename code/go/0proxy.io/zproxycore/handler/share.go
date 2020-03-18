package handler

import (
	"context"
	"net/http"
	"path/filepath"

	"0proxy.io/zproxycore/common"
	"github.com/0chain/gosdk/zboxcore/fileref"
	"github.com/0chain/gosdk/zboxcore/sdk"
)

func Share(ctx context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodGet {
		return nil, common.NewError("invalid_method", "Invalid method for delete endpoint, Use GET.")
	}
	allocation := r.FormValue("allocation")
	clientJSON := r.FormValue("client_json")
	remotePath := r.FormValue("remote_path")
	refereeClientID := r.FormValue("referee_client_id")
	encryptionpublickey := r.FormValue(("encryption_public_key"))

	err := initSDK(clientJSON)
	if err != nil {
		return nil, common.NewError("sdk_not_initialized", "Unable to initialize gosdk with the given client details")
	}

	allocationObj, err := sdk.GetAllocation(allocation)
	if err != nil {
		return nil, err
	}
	refType := fileref.FILE
	statsMap, err := allocationObj.GetFileStats(remotePath)
	if err != nil {
		return nil, err
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

	ref, err := allocationObj.GetAuthTicket(remotePath, fileName, refType, refereeClientID, encryptionpublickey)
	if err != nil {
		return nil, err
	}

	var response struct {
		authTicket string `json:"auth_ticket"`
	}
	response.authTicket = ref
	return response, nil
}
