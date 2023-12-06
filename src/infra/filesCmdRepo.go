package infra

import (
	"errors"
	"log"
	"os"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type FilesCmdRepo struct {
}

func (repo FilesCmdRepo) Add(addUnixFile dto.AddUnixFile) error {
	if !addUnixFile.Type.IsDir() {
		_, err := os.Create(addUnixFile.Path.String())
		if err != nil {
			log.Printf("CreateUnixFileError: %s", err)
			return errors.New("CreateUnixFileError")
		}

		return repo.UpdatePermissions(addUnixFile.Path, addUnixFile.Permissions)
	}

	err := os.MkdirAll(addUnixFile.Path.String(), addUnixFile.Permissions.GetFileMode())
	if err != nil {
		log.Printf("CreateUnixDirectoryError: %s", err)
		return errors.New("CreateUnixDirectoryError")
	}

	return nil
}

func (repo FilesCmdRepo) UpdatePermissions(
	unixFilePath valueObject.UnixFilePath,
	unixFilePermissions valueObject.UnixFilePermissions,
) error {
	err := os.Chmod(unixFilePath.String(), unixFilePermissions.GetFileMode())
	if err != nil {
		log.Printf("ChmodUnixFileError: %s", err)
		return errors.New("ChmodUnixFileError")
	}

	return nil
}
