package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
)

const FilesRepo = "files/"

func writeFile(file multipart.File, filePath string) (string, error) {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	return f.Name(), err
}

func deleletFile(filePath string) error {
	return os.RemoveAll(filePath)
}

func readFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func getPath(allocation, fileName string) string {
	return FilesRepo + allocation + "/" + fileName
}

func getPathForStream(allocation, fileName string, start, end int) string {
	return FilesRepo + allocation + "/" + fmt.Sprintf("%d-%d-%s", start, end, fileName)
}

func createDirIfNotExists(allocation string) {
	allocationDir := FilesRepo + allocation
	if _, err := os.Stat(allocationDir); os.IsNotExist(err) {
		os.Mkdir(allocationDir, 0777)
	}
}
