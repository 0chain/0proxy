package handler

import (
	"context"
	"net/http"

	"0proxy.io/core/common"
	"github.com/0chain/gosdk/zboxcore/sdk"
)

func GetPublicEncryptionKey(ctx context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodGet {
		return nil, methodError("publiceEcryptionKey", http.MethodGet)
	}

	clientJSON := r.FormValue("client_json")
	if len(clientJSON) == 0 {
		return nil, common.NewError("invalid_param", "Please provide client_json for the client")
	}

	err := initSDK(clientJSON)
	if err != nil {
		return nil, common.NewError("sdk_not_initialized", "Unable to initialize gosdk with the given client details")
	}

	key, err := sdk.GetClientEncryptedPublicKey()
	if err != nil {
		return nil, common.NewError("get_public_encryption_key_failed", err.Error())
	}

	resp := &struct {
		Key string `json:"key"`
	}{Key: key}

	return resp, nil
}
