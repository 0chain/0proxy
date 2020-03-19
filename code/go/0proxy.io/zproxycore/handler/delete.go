package handler

import (
	"context"
	"net/http"

	"0proxy.io/core/common"
	"github.com/0chain/gosdk/zboxcore/sdk"
)

// Delete is to delete a file in dStorage
func Delete(ctx context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodDelete {
		return nil, methodError("delete", http.MethodDelete)
	}

	allocation := r.FormValue("allocation")
	clientJSON := r.FormValue("client_json")
	err := validateClientDetails(allocation, clientJSON)
	if err != nil {
		return nil, err
	}

	remotePath := r.FormValue("remote_path")
	if len(remotePath) == 0 {
		return nil, common.NewError("invalid_param", "Please provide remote_path for delete")
	}

	err = initSDK(clientJSON)
	if err != nil {
		return nil, common.NewError("sdk_not_initialized", "Unable to initialize gosdk with the given client details")
	}

	allocationObj, err := sdk.GetAllocation(allocation)
	if err != nil {
		return nil, common.NewError("get_allocation_failed", err.Error())
	}

	err = allocationObj.DeleteFile(remotePath)
	if err != nil {
		return nil, common.NewError("delete_object_failed", err.Error())
	}

	resp := response{
		Message: "Delete done successfully",
	}
	return resp, nil
}
