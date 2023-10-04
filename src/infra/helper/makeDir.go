package infraHelper

import (
	"os"
	"regexp"
)

func MakeDir(dirPath string) error {
	_, err := os.Stat(dirPath)
	if err != nil {
		matchErr, _ := regexp.MatchString("no such file or directory", err.Error())
		if !matchErr {
			return err
		}
	}

	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
