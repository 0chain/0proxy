package handler

import (
	"context"
	"net/http"

	"0proxy.io/zproxycore/common"
	"github.com/0chain/gosdk/zboxcore/sdk"
)

func Delete(ctx context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodDelete {
		return nil, common.NewError("invalid_method", "Invalid method for delete endpoint, Use DELETE.")
	}
	allocation := r.FormValue("allocation")
	clientJSON := r.FormValue("client_json")
	remotePath := r.FormValue("remote_path")

	err := initSDK(clientJSON)
	if err != nil {
		return nil, common.NewError("sdk_not_initialized", "Unable to initialize gosdk with the given client details")
	}

	allocationObj, err := sdk.GetAllocation(allocation)
	if err != nil {
		return nil, err
	}

	err = allocationObj.DeleteFile(remotePath)
	if err != nil {
		return nil, err
	}

	var response struct {
		msg string
	}
	response.msg = "Delete done"
	return response, nil
}
