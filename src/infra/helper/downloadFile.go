package infraHelper

import (
	"errors"
	"log"

	tkInfra "github.com/goinfinite/tk/src/infra"
)

func DownloadFile(url string, filePath string) error {
	_, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "wget",
		Args:    []string{"-q", "--no-check-certificate", "-O", filePath, url},
	}).Run()
	if err != nil {
		log.Printf("DownloadFileError: %s", err)
		return errors.New("DownloadFileError")
	}

	return nil
}
