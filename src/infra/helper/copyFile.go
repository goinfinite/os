package infraHelper

import (
	"errors"
	"io"
	"os"
)

func CopyFile(srcPath string, dstPath string) error {
	srcFile, err := os.OpenFile(srcPath, os.O_RDWR, 0644)
	if err != nil {
		return errors.New("OpenSourceFilePathError: " + err.Error())
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.New("OpenDestinationFilePathError: " + err.Error())
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return errors.New("CopyFileError: " + err.Error())
	}

	return nil
}
