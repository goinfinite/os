package useCase

import (
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
		return err
	}

	log.Printf(
		"File '%s' content updated.",
		updateUnixFileContent.Path.GetFileName().String(),
	)

	return nil
}
