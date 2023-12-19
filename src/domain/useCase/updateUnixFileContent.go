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
	filePath := updateUnixFileContent.Path

	unixFileExists, err := filesQueryRepo.Exists(filePath)
	if err != nil {
		return err
	}

	if unixFileExists {
		return errors.New("FileDoesNotExists")
	}

	isDir, err := filesQueryRepo.IsDir(filePath)
	if err != nil {
		return err
	}

	if isDir {
		return errors.New("FilePathIsDir")
	}

	err = filesCmdRepo.UpdateContent(updateUnixFileContent)
	if err != nil {
		return err
	}

	fileName, _ := updateUnixFileContent.Path.GetFileName()
	log.Printf("File '%s' content updated.", fileName.String())

	return nil
}
