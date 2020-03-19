package handler

import (
	"context"
	"net/http"

	"0proxy.io/core/common"
	"github.com/0chain/gosdk/zboxcore/sdk"
)

// Move is to move file in dStorage
func Move(ctx context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodPut {
		return nil, methodError("move", http.MethodPut)
	}
	allocation := r.FormValue("allocation")
	clientJSON := r.FormValue("client_json")
	err := validateClientDetails(allocation, clientJSON)
	if err != nil {
		return nil, err
	}

	remotePath := r.FormValue("remote_path")
	destPath := r.FormValue("dest_path")
	if len(remotePath) == 0 || len(destPath) == 0 {
		return nil, common.NewError("invalid_param", "Please provide remote_path and dest_path for move")
	}

	err = initSDK(clientJSON)
	if err != nil {
		return nil, common.NewError("sdk_not_initialized", "Unable to initialize gosdk with the given client details")
	}

	allocationObj, err := sdk.GetAllocation(allocation)
	if err != nil {
		return nil, common.NewError("get_allocation_failed", err.Error())
	}

	err = allocationObj.MoveObject(remotePath, destPath)
	if err != nil {
		return nil, common.NewError("move_object_failed", err.Error())
	}

	resp := response{
		Message: "Move done successfully",
	}
	return resp, nil
}
