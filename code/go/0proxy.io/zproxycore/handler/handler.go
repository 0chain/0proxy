package handler

import (
	"fmt"
	"os"
	"sync"

	"0proxy.io/core/config"
	"github.com/0chain/gosdk/zboxcore/sdk"
	"gopkg.in/cheggaaa/pb.v1"
)

func initSDK(clientJSON string) error {
	return sdk.InitStorageSDK(clientJSON,
		config.Configuration.Miners,
		config.Configuration.Sharders,
		config.Configuration.ChainID,
		config.Configuration.SignatureScheme,
		nil)
}

func (s *StatusBar) Started(allocationId, filePath string, op int, totalBytes int) {
	s.b = pb.StartNew(totalBytes)
	s.b.Set(0)
}
func (s *StatusBar) InProgress(allocationId, filePath string, op int, completedBytes int) {
	s.b.Set(completedBytes)
}

func (s *StatusBar) Completed(allocationId, filePath string, filename string, mimetype string, size int, op int) {
	if s.b != nil {
		s.b.Finish()
	}
	s.success = true
	defer s.wg.Done()
	fmt.Println("Status completed callback. Type = " + mimetype + ". Name = " + filename)
}

func (s *StatusBar) Error(allocationID string, filePath string, op int, err error) {
	if s.b != nil {
		s.b.Finish()
	}
	s.success = false
	defer s.wg.Done()
	PrintError("Error in file operation." + err.Error())
}

func PrintError(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

type StatusBar struct {
	b       *pb.ProgressBar
	wg      *sync.WaitGroup
	success bool
}
