package common

import (
	"context"
	"net/http"

	"0proxy.io/core/common"
	"0proxy.io/zproxycore/handler"

	"github.com/gorilla/mux"
)

// SetupHandlers is to setup handlers
func SetupHandlers(r *mux.Router) {
	r.HandleFunc("/upload", common.UserRateLimit(common.ToJSONResponse(Upload)))
	r.HandleFunc("/download", common.UserRateLimit(common.ToFileResponse(Download)))
	r.HandleFunc("/delete", common.UserRateLimit(common.ToJSONResponse(Delete)))
	r.HandleFunc("/share", common.UserRateLimit(common.ToJSONResponse(Share)))
	r.HandleFunc("/copy", common.UserRateLimit(common.ToJSONResponse(Copy)))
	r.HandleFunc("/rename", common.UserRateLimit(common.ToJSONResponse(Rename)))
	r.HandleFunc("/move", common.UserRateLimit(common.ToJSONResponse(Move)))
}

// Upload is for file upload
func Upload(ctx context.Context, r *http.Request) (interface{}, error) {
	return handler.Upload(ctx, r)
}

// Move is to move a file
func Move(ctx context.Context, r *http.Request) (interface{}, error) {
	return handler.Move(ctx, r)
}

// Delete is for file delete
func Delete(ctx context.Context, r *http.Request) (interface{}, error) {
	return handler.Delete(ctx, r)
}

// Share to share a file and get authTicket
func Share(ctx context.Context, r *http.Request) (interface{}, error) {
	return handler.Share(ctx, r)
}

// Rename is to rename file
func Rename(ctx context.Context, r *http.Request) (interface{}, error) {
	return handler.Rename(ctx, r)
}

// Copy is to copy a file
func Copy(ctx context.Context, r *http.Request) (interface{}, error) {
	return handler.Copy(ctx, r)
}

// Download is to download a file
func Download(w http.ResponseWriter, r *http.Request) {
	handler.Download(w, r)
}
