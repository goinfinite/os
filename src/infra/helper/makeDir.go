package infraHelper

import (
	"os"
)

func MakeDir(dirPath string) error {
	_, err := os.Stat(dirPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
