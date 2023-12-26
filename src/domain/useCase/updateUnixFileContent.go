package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func UpdateUnixFileContent(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	updateUnixFileContent dto.UpdateUnixFileContent,
) error {
	err := filesCmdRepo.UpdateContent(updateUnixFileContent)
	if err != nil {
		log.Printf("UpdateUnixFileContentInfraError: %s", err.Error())
		return errors.New("UpdateUnixFileContentInfraError")
	}

	log.Printf(
		"File '%s' content updated.",
		updateUnixFileContent.SourcePath.GetFileName().String(),
	)

	return nil
}
