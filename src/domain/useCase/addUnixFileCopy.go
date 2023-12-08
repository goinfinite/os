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

	unixFiles, err := filesQueryRepo.Get(filePath)

	if err != nil || len(unixFiles) < 1 {
		return errors.New("PathDoesNotExists")
	}

	fileType := "File"
	fileIsDir, err := filePath.IsDir()
	if err != nil {
		log.Printf("PathIsDirError: %s", err)
		return errors.New("PathIsDirError")
	}

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
