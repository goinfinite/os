package useCase

import (
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func AddUnixFile(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	addUnixFile dto.AddUnixFile,
) error {
	err := filesCmdRepo.Create(addUnixFile)
	if err != nil {
		return err
	}

	log.Printf(
		"File '%s' created to '%s'.",
		addUnixFile.Path.GetFileName().String(),
		addUnixFile.Path.GetFileDir().String(),
	)

	return nil
}
