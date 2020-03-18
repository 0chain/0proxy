package common

import (
	"context"
	"net/http"

	"0proxy.io/core/common"
	. "0proxy.io/core/logging"
	"0proxy.io/zproxycore/handler"

	"github.com/gorilla/mux"
)

func SetupHandlers(r *mux.Router) {
	r.HandleFunc("/upload", common.ToJSONResponse(Upload))
	r.HandleFunc("/download", common.ToFileResponse(Download))
	r.HandleFunc("/delete", common.ToJSONResponse(Delete))
	r.HandleFunc("/share", common.ToJSONResponse(Share))
	r.HandleFunc("/copy", common.ToJSONResponse(Copy))
	r.HandleFunc("/rename", common.ToJSONResponse(Rename))
}

func Upload(ctx context.Context, r *http.Request) (interface{}, error) {
	return handler.Upload(ctx, r)
}

func Delete(ctx context.Context, r *http.Request) (interface{}, error) {
	return handler.Delete(ctx, r)
}
func Share(ctx context.Context, r *http.Request) (interface{}, error) {
	return handler.Share(ctx, r)
}
func Rename(ctx context.Context, r *http.Request) (interface{}, error) {
	return handler.Rename(ctx, r)
}
func Copy(ctx context.Context, r *http.Request) (interface{}, error) {
	return handler.Copy(ctx, r)
}

func Download(w http.ResponseWriter, r *http.Request) {
	handler.Download(w, r)
}
