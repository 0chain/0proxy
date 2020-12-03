package worker

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"0proxy.io/core/config"
	. "0proxy.io/core/logging"
	"0proxy.io/zproxycore/handler"
	"go.uber.org/zap"
)

func SetupWorkers(ctx context.Context) {
	go CacheCleanUp(ctx)
}

func CacheCleanUp(ctx context.Context) {
	var iterInprogress = false
	ticker := time.NewTicker(time.Duration(config.Configuration.CleanUpWorkerMinutes) * time.Minute)
	for true {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !iterInprogress {
				iterInprogress = true
				now := time.Now()
				err := filepath.Walk(handler.FilesRepo,
					func(path string, fi os.FileInfo, err error) error {
						if fi.IsDir() {
							return nil
						}
						if diff := now.Sub(fi.ModTime()); diff > 10*time.Minute {
							os.RemoveAll(path)
						}
						return nil
					})
				if err != nil {
					Logger.Error("Clean up worker failed to walk through files directory", zap.Error(err))
					continue
				}
				iterInprogress = false
				Logger.Info("Clean up worker cycle completed.")
			}
		}
	}
}
