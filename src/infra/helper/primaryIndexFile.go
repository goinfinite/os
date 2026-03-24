package infraHelper

import (
	"errors"
	"os"

	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

var fileClerk = tkInfra.FileClerk{}

var (
	IndexFileTemplatePath string = infraEnvs.VirtualHostsConfDir + "/index.html"
	IndexFilePath         string = infraEnvs.PrimaryPublicDir + "/index.html"
	IndexFileBackupPath   string = infraEnvs.PrimaryPublicDir + "/../index.html.backup"
)

func BackupPrimaryIndexFile() error {
	if !fileClerk.FileExists(IndexFilePath) {
		return nil
	}

	if fileClerk.FileExists(IndexFileBackupPath) {
		err := os.Remove(IndexFileBackupPath)
		if err != nil {
			return errors.New(
				"RemoveIndexBackupFileError: " + err.Error(),
			)
		}
	}

	err := os.Rename(IndexFilePath, IndexFileBackupPath)
	if err != nil {
		return errors.New("MoveIndexFileError: " + err.Error())
	}

	err = UpdateOwnershipForWebServerUse(
		IndexFileBackupPath, false, false,
	)
	if err != nil {
		return errors.New(
			"UpdateOwnershipForWebServerUseError: " + err.Error(),
		)
	}

	return nil
}

func RestorePrimaryIndexFile() error {
	if fileClerk.FileExists(IndexFilePath) {
		return nil
	}

	restorableIndexFilePath := IndexFileTemplatePath
	if fileClerk.FileExists(IndexFileBackupPath) {
		restorableIndexFilePath = IndexFileBackupPath
	}

	err := fileClerk.CopyFile(restorableIndexFilePath, IndexFilePath)
	if err != nil {
		return errors.New("CopyIndexFileError: " + err.Error())
	}

	err = UpdateOwnershipForWebServerUse(IndexFilePath, false, false)
	if err != nil {
		return errors.New(
			"UpdateOwnershipForWebServerUseError: " + err.Error(),
		)
	}

	if fileClerk.FileExists(IndexFileBackupPath) {
		err = os.Remove(IndexFileBackupPath)
		if err != nil {
			return errors.New(
				"RemoveIndexBackupFileError: " + err.Error(),
			)
		}
	}

	return nil
}
