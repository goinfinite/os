package infraHelper

import (
	"errors"
	"os"
)

func GetFileContent(filePath string) (string, error) {
	fileExists := FileExists(filePath)
	if !fileExists {
		return "", errors.New("FileNotFound")
	}

	_, err := os.Stat(filePath)
	if err != nil {
		return "", errors.New("FailedToGetFileInfo: " + err.Error())
	}

	fileContentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", errors.New("FailedToGetFileContent: " + err.Error())
	}
	fileContentStr := string(fileContentBytes)

	return fileContentStr, nil
}
