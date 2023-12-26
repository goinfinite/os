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
	filePath := updateUnixFile.Path

	fileName, _ := filePath.GetFileName()
	fileDir, _ := filePath.GetFileDir()

	filePermissions := updateUnixFile.Permissions
	if filePermissions != nil {
		err := filesCmdRepo.UpdatePermissions(filePath, *filePermissions)
		if err != nil {
			log.Printf("UpdateFilePermissionsError: %s", err.Error())
			return errors.New("UpdateFilePermissionsError")
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
		log.Printf("MoveFileError: %s", err.Error())
		return errors.New("MoveFileError")
	}

	fileDestinationDir, _ := updateUnixFile.DestinationPath.GetFileDir()
	log.Printf(
		"File '%s' moved from %s to '%s'.",
		fileName.String(),
		fileDir.String(),
		fileDestinationDir.String(),
	)

	return nil
}
