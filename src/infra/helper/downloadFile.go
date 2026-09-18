package infraHelper

import (
	"errors"
	"log/slog"

	tkInfra "github.com/goinfinite/tk/src/infra"
)

func DownloadFile(url string, filePath string) error {
	_, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "wget",
		Args:    []string{"-q", "--no-check-certificate", "-O", filePath, url},
	}).Run()
	if err != nil {
		slog.Error("DownloadFileError", slog.String("err", err.Error()))
		return errors.New("DownloadFileError")
	}

	return nil
}
