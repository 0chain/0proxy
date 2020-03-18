package handler

import (
	"context"
	"net/http"

	"0proxy.io/zproxycore/common"
	"github.com/0chain/gosdk/zboxcore/sdk"
)

func Move(ctx context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodPut {
		return nil, common.NewError("invalid_method", "Invalid method for delete endpoint, Use PUT.")
	}
	allocation := r.FormValue("allocation")
	clientJSON := r.FormValue("client_json")
	remotePath := r.FormValue("remote_path")
	destPath := r.FormValue("dest_path")

	err := initSDK(clientJSON)
	if err != nil {
		return nil, common.NewError("sdk_not_initialized", "Unable to initialize gosdk with the given client details")
	}

	allocationObj, err := sdk.GetAllocation(allocation)
	if err != nil {
		return nil, err
	}

	err = allocationObj.MoveObject(remotePath, destPath)
	if err != nil {
		return nil, err
	}

	var response struct {
		msg string
	}
	response.msg = "Move done"
	return response, nil
}
