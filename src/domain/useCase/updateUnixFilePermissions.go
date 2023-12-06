package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func UpdateUnixFilePermissions(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	unixFilePath valueObject.UnixFilePath,
	unixFilePermissions valueObject.UnixFilePermissions,
	unixFileType valueObject.UnixFileType,
) error {
	unixFiles, err := filesQueryRepo.Get(unixFilePath)

	if err != nil && len(unixFiles) < 1 {
		return errors.New("FileDoesNotExists")
	}

	err = filesCmdRepo.UpdatePermissions(
		unixFilePath,
		unixFilePermissions,
		unixFileType,
	)
	if err != nil {
		return errors.New("UpdateFilePermissionsError")
	}

	fileName, _ := unixFilePath.GetFileName()
	log.Printf("File '%s' permissions updated.", fileName.String())

	return nil
}
