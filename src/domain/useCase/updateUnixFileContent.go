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

	unixFiles, err := filesQueryRepo.Get(filePath)
	if err != nil || len(unixFiles) < 1 {
		return errors.New("PathDoesNotExists")
	}

	isDir, err := filePath.IsDir()
	if err != nil {
		log.Printf("PathIsDirError: %s", err)
		return errors.New("PathIsDirError")
	}

	if isDir {
		return errors.New("FilePathIsDir")
	}

	err = filesCmdRepo.UpdateContent(updateUnixFileContent)
	if err != nil {
		return errors.New("UpdateFileContentError")
	}

	fileName, _ := updateUnixFileContent.Path.GetFileName()
	log.Printf("File '%s' content updated.", fileName.String())

	return nil
}
