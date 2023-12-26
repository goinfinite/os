package useCase

import (
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
		log.Printf("CreateFileError: %s", err.Error())
		return err
	}

	log.Printf(
		"File '%s' created in '%s'.",
		createUnixFile.SourcePath.GetFileName().String(),
		createUnixFile.SourcePath.GetFileDir().String(),
	)

	return nil
}
