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
	fileType := addUnixFile.Type.GetWithFirstLetterUpperCase()

	unixFileExists := filesQueryRepo.Exists(addUnixFile.Path)
	if unixFileExists {
		return errors.New(fileType + "AlreadyExists")
	}

	err := filesCmdRepo.Add(addUnixFile)
	if err != nil {
		return errors.New("Add" + fileType + "Error")
	}

	fileName, _ := addUnixFile.Path.GetFileName()
	fileDir, _ := addUnixFile.Path.GetFileDir()
	log.Printf("%s '%s' added to '%s'.", fileType, fileName.String(), fileDir.String())

	return nil
}
