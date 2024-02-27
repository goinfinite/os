package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func CreateUnixFile(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	createUnixFile dto.CreateUnixFile,
) error {
	err := filesCmdRepo.Create(createUnixFile)
	if err != nil {
		log.Printf("CreateUnixFileInfraError: %s", err.Error())
		return errors.New("CreateUnixFileInfraError")
	}

	log.Printf(
		"File '%s' created in '%s'.",
		createUnixFile.FilePath.GetFileName().String(),
		createUnixFile.FilePath.GetFileDir().String(),
	)

	return nil
}
