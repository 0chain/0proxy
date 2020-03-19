package handler

import (
	"fmt"
	"os"
	"sync"

	"0proxy.io/core/common"
	"0proxy.io/core/config"
	"github.com/0chain/gosdk/zboxcore/sdk"
	"gopkg.in/cheggaaa/pb.v1"
)

type response struct {
	Message string `json:"message"`
}

func methodError(endpoint, method string) error {
	return common.NewError("invalid_method", fmt.Sprintf("Invalid method for %s endpoint, Use %s", endpoint, method))
}

func initSDK(clientJSON string) error {
	return sdk.InitStorageSDK(clientJSON,
		config.Configuration.Miners,
		config.Configuration.Sharders,
		config.Configuration.ChainID,
		config.Configuration.SignatureScheme,
		nil)
}

func validateClientDetails(allocation, clientJSON string) error {
	if len(allocation) == 0 || len(clientJSON) == 0 {
		return common.NewError("invalid_param", "Please provide allocation and client_json for the client")
	}
	return nil
}

// Started for statusBar
func (s *StatusBar) Started(allocationID, filePath string, op int, totalBytes int) {
	s.b = pb.StartNew(totalBytes)
	s.b.Set(0)
}

// InProgress for statusBar
func (s *StatusBar) InProgress(allocationID, filePath string, op int, completedBytes int) {
	s.b.Set(completedBytes)
}

// Completed for statusBar
func (s *StatusBar) Completed(allocationID, filePath string, filename string, mimetype string, size int, op int) {
	if s.b != nil {
		s.b.Finish()
	}
	s.success = true
	defer s.wg.Done()
	fmt.Println("Status completed callback. Type = " + mimetype + ". Name = " + filename)
}

// Error for statusBar
func (s *StatusBar) Error(allocationID string, filePath string, op int, err error) {
	if s.b != nil {
		s.b.Finish()
	}
	s.success = false
	defer s.wg.Done()
	PrintError("Error in file operation." + err.Error())
}

// PrintError is to print error
func PrintError(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

// StatusBar is to check status of any operation
type StatusBar struct {
	b       *pb.ProgressBar
	wg      *sync.WaitGroup
	success bool
}
