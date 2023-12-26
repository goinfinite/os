package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func UpdateUnixFile(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	updateUnixFile dto.UpdateUnixFile,
) error {
	filePath := updateUnixFile.SourcePath

	fileName := filePath.GetFileName()
	fileDir := filePath.GetFileDir()

	filePermissions := updateUnixFile.Permissions
	if filePermissions != nil {
		err := filesCmdRepo.UpdatePermissions(filePath, *filePermissions)
		if err != nil {
			log.Printf("UpdateUnixFilePermissionsInfraError: %s", err.Error())
			return errors.New("UpdateUnixFilePermissionsInfraError")
		}

		log.Printf(
			"File '%s' (%s) permissions updated to '%s'.",
			fileName.String(),
			fileDir.String(),
			filePermissions.String(),
		)
	}

	if updateUnixFile.DestinationPath == nil {
		return nil
	}

	err := filesCmdRepo.Move(updateUnixFile)
	if err != nil {
		log.Printf("MoveUnixFileInfraError: %s", err.Error())
		return errors.New("MoveUnixFileInfraError")
	}

	log.Printf(
		"File '%s' moved from %s to '%s'.",
		fileName.String(),
		fileDir.String(),
		updateUnixFile.DestinationPath.GetFileDir().String(),
	)

	return nil
}
