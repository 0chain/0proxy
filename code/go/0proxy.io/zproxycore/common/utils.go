package common

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
)

func WriteFile(file multipart.File, filePath string) (string, error) {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	return f.Name(), err
}

func DeleletFile(filePath string) error {
	return os.RemoveAll(filePath)
}

func ReadFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func GetPath(allocation, fileName string) string {
	return "./" + allocation + "/" + fileName
}

func CreateDirIfNotExists(allocation string) {
	allocationDir := "./" + allocation
	if _, err := os.Stat(allocationDir); os.IsNotExist(err) {
		os.Mkdir(allocationDir, 0777)
	}
}
