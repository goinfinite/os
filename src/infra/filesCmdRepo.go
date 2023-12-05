package infra

import (
	"errors"
	"log"
	"os"

	"github.com/speedianet/os/src/domain/dto"
)

type FilesCmdRepo struct {
}

func (repo FilesCmdRepo) Add(addUnixFile dto.AddUnixFile) error {
	filePathStr := addUnixFile.Path.String()

	_, err := os.Create(filePathStr)
	if err != nil {
		log.Printf("CreateUnixFileError: %s", err)
		return errors.New("CreateUnixFileError")
	}

	err = os.Chmod(filePathStr, addUnixFile.Permissions.GetFileMode())
	if err != nil {
		log.Printf("ChmodUnixFileError: %s", err)
		return errors.New("ChmodUnixFileError")
	}

	return nil
}
