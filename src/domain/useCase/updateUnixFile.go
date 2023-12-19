package useCase

import (
	"errors"
	"fmt"
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

	unixFileExists, err := filesQueryRepo.Exists(filePath)
	if err != nil {
		return err
	}

	if unixFileExists {
		return errors.New("PathDoesNotExists")
	}

	fileType := "File"
	fileIsDir, err := filesQueryRepo.IsDir(filePath)
	if err != nil {
		return err
	}

	if fileIsDir {
		fileType = "Dir"
	}

	fileName, _ := filePath.GetFileName()
	fileDir, _ := filePath.GetFileDir()

	if updateUnixFile.Permissions != nil {
		filePermissions := *updateUnixFile.Permissions

		err = filesCmdRepo.UpdatePermissions(filePath, filePermissions)
		if err != nil {
			return errors.New("Update" + fileType + "PermissionsError")
		}

		log.Printf(
			"%s '%s' (%s) permissions updated to '%s'.",
			fileType,
			fileName.String(),
			fileDir.String(),
			filePermissions.String(),
		)
	}

	if updateUnixFile.DestinationPath != nil {
		fileDestinationPath := *updateUnixFile.DestinationPath

		fileDestinationDir, _ := fileDestinationPath.GetFileDir()

		toBeRenamed := fileDir.String() == fileDestinationDir.String()

		err = filesCmdRepo.Move(updateUnixFile.Path, fileDestinationPath)
		if err != nil {
			processToBeExecuted := "Move"
			if toBeRenamed {
				processToBeExecuted = "Rename"
			}

			failureMessage := fmt.Sprintf("%s%sError", processToBeExecuted, fileType)
			return errors.New(failureMessage)
		}

		successMessage := fmt.Sprintf(
			"%s '%s' moved from %s to '%s'.",
			fileType,
			fileName.String(),
			fileDir.String(),
			fileDestinationDir.String(),
		)

		if toBeRenamed {
			fileDestinationName, _ := fileDestinationPath.GetFileName()
			successMessage = fmt.Sprintf(
				"%s '%s' renamed to '%s'.",
				fileType,
				fileName.String(),
				fileDestinationName.String(),
			)
		}
		log.Printf(successMessage)
	}

	return nil
}
