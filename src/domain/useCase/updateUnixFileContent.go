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
	unixFiles, err := filesQueryRepo.Get(updateUnixFileContent.Path)

	if err != nil && len(unixFiles) < 1 {
		return errors.New("FileDoesNotExists")
	}

	err = filesCmdRepo.UpdateContent(updateUnixFileContent)
	if err != nil {
		return errors.New("UpdateFileContentError")
	}

	fileName, _ := updateUnixFileContent.Path.GetFileName()
	log.Printf("File '%s' content updated.", fileName.String())

	return nil
}
