package infraHelper

import (
	"errors"
	"log"
	"os"
)

func ReadFile(filePath string) (string, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		log.Printf("FailedToReadFile: %v", err)
		return "", errors.New("FailedToReadFile")
	}

	fileContentBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("FailedToReadFile: %v", err)
		return "", errors.New("FailedToReadFile")
	}
	fileContentStr := string(fileContentBytes)

	return fileContentStr, nil
}
