package handler

import (
	"context"
	"net/http"

	"0proxy.io/core/common"
	"github.com/0chain/gosdk/zboxcore/sdk"
)

// Rename is to rename file in dStorage
func Rename(ctx context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodPut {
		return nil, methodError("rename", http.MethodPut)
	}
	allocation := r.FormValue("allocation")
	clientJSON := r.FormValue("client_json")
	err := validateClientDetails(allocation, clientJSON)
	if err != nil {
		return nil, err
	}

	remotePath := r.FormValue("remote_path")
	newName := r.FormValue("new_name")
	if len(remotePath) == 0 || len(newName) == 0 {
		return nil, common.NewError("invalid_param", "Please provide remote_path and new_name for rename")
	}

	err = initSDK(clientJSON)
	if err != nil {
		return nil, common.NewError("sdk_not_initialized", "Unable to initialize gosdk with the given client details")
	}

	allocationObj, err := sdk.GetAllocation(allocation)
	if err != nil {
		return nil, common.NewError("get_allocation_failed", err.Error())
	}

	err = allocationObj.RenameObject(remotePath, newName)
	if err != nil {
		return nil, common.NewError("rename_object_failed", err.Error())
	}

	resp := response{
		Message: "Rename done successfully",
	}
	return resp, nil
}
