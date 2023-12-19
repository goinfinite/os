package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func AddUnixFileCopy(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	addUnixFileCopy dto.AddUnixFileCopy,
) error {
	filePath := addUnixFileCopy.OriginPath

	fileIsDir, err := filesQueryRepo.IsDir(filePath)
	if err != nil {
		return errors.New("PathIsDirError")
	}

	fileExists, err := filesQueryRepo.Exists(filePath)
	if err != nil {
		return errors.New("PathExistsError")
	}

	if fileExists {
		return errors.New("PathAlreadyExists")
	}

	fileType := "File"
	if fileIsDir {
		fileType = "Dir"
	}

	err = filesCmdRepo.Copy(addUnixFileCopy)
	if err != nil {
		return errors.New("Add" + fileType + "CopyError")
	}

	fileName, _ := addUnixFileCopy.OriginPath.GetFileName()
	fileDir, _ := addUnixFileCopy.OriginPath.GetFileDir()

	fileDestinationPath := addUnixFileCopy.DestinationPath
	fileDestinationName, _ := fileDestinationPath.GetFileName()
	fileDestinationDir, _ := fileDestinationPath.GetFileDir()
	log.Printf(
		"%s '%s' (%s) copy added to '%s' with name '%s'.",
		fileType,
		fileName.String(),
		fileDir.String(),
		fileDestinationDir.String(),
		fileDestinationName.String(),
	)

	return nil
}
