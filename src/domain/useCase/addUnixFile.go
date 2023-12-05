package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func AddUnixFile(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	addUnixFile dto.AddUnixFile,
) error {
	unixFiles, err := filesQueryRepo.Get(addUnixFile.Path)
	if err != nil {
		return err
	}

	if len(unixFiles) > 0 {
		return errors.New("FileAlreadyExists")
	}

	err = filesCmdRepo.Add(addUnixFile)
	if err != nil {
		return errors.New("AddFileError")
	}

	log.Printf("File '%s' added to '%s'.", addUnixFile.Name.String(), addUnixFile.Path)

	return nil
}
