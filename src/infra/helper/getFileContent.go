package infraHelper

import (
	"errors"
	"log"
	"os"
)

func GetFileContent(filePath string) (string, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		log.Printf("FailedToGetFileContent: %v", err)
		return "", errors.New("FailedToGetFileContent")
	}

	fileContentBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("FailedToGetFileContent: %v", err)
		return "", errors.New("FailedToGetFileContent")
	}
	fileContentStr := string(fileContentBytes)

	return fileContentStr, nil
}
