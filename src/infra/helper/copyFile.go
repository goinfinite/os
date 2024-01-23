package infraHelper

import (
	"errors"
	"io"
	"log"
	"os"
)

func CopyFile(srcPath string, dstPath string) error {
	fileFlags := os.O_RDWR | os.O_CREATE | os.O_TRUNC

	srcFile, err := os.OpenFile(dstPath, fileFlags, 0644)
	if err != nil {
		return errors.New("OpenSourceFileError: " + err.Error())
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dstPath, fileFlags, 0644)
	if err != nil {
		return errors.New("OpenDestinationFileError: " + err.Error())
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Printf("CopyFileError: %s", err)
		return errors.New("CopyFileError")
	}

	return nil
}
