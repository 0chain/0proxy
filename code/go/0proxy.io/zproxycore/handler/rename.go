package handler

import (
	"context"
	"net/http"

	"0proxy.io/zproxycore/common"
	"github.com/0chain/gosdk/zboxcore/sdk"
)

func Rename(ctx context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodPut {
		return nil, common.NewError("invalid_method", "Invalid method for delete endpoint, Use PUT.")
	}
	allocation := r.FormValue("allocation")
	clientJSON := r.FormValue("client_json")
	remotePath := r.FormValue("remote_path")
	newName := r.FormValue("new_name")

	err := initSDK(clientJSON)
	if err != nil {
		return nil, common.NewError("sdk_not_initialized", "Unable to initialize gosdk with the given client details")
	}

	allocationObj, err := sdk.GetAllocation(allocation)
	if err != nil {
		return nil, err
	}

	err = allocationObj.RenameObject(remotePath, newName)
	if err != nil {
		return nil, err
	}

	var response struct {
		msg string
	}
	response.msg = "Rename done"
	return response, nil
}
