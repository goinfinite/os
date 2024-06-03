package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

type UpdateUnixFiles struct {
	filesCmdRepo repository.FilesCmdRepo
}

func NewUpdateUnixFiles(
	filesCmdRepo repository.FilesCmdRepo,
) UpdateUnixFiles {
	return UpdateUnixFiles{
		filesCmdRepo: filesCmdRepo,
	}
}

func (uc UpdateUnixFiles) updateFailureFactory(
	filePath valueObject.UnixFilePath,
	errMessage string,
) (valueObject.UpdateProcessFailure, error) {
	var updateProcessFailure valueObject.UpdateProcessFailure

	failureReason, err := valueObject.NewFailureReason(errMessage)
	if err != nil {
		return updateProcessFailure, err
	}

	return valueObject.NewUpdateProcessFailure(
		filePath,
		failureReason,
	), nil
}

func (uc UpdateUnixFiles) updateFilePermissions(
	sourcePath valueObject.UnixFilePath,
	permissions valueObject.UnixFilePermissions,
) error {
	err := uc.filesCmdRepo.UpdatePermissions(
		sourcePath,
		permissions,
	)
	if err != nil {
		log.Printf("UpdateFilePermissionsError: %s", err.Error())
		return errors.New("UpdateFilePermissionsInfraError")
	}

	log.Printf(
		"File '%s' (%s) permissions updated to '%s'.",
		sourcePath.GetFileName().String(),
		sourcePath.GetFileDir().String(),
		permissions.String(),
	)

	return nil
}

func (uc UpdateUnixFiles) moveFile(
	sourcePath valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
) error {
	shouldOverwrite := false
	err := uc.filesCmdRepo.Move(
		sourcePath,
		destinationPath,
		shouldOverwrite,
	)
	if err != nil {
		log.Printf("MoveFileError: %s", err.Error())
		return errors.New("MoveFileInfraError")
	}

	return nil
}

func (uc UpdateUnixFiles) updateFileContent(
	sourcePath valueObject.UnixFilePath,
	encodedContent valueObject.EncodedContent,
) error {
	err := uc.filesCmdRepo.UpdateContent(sourcePath, encodedContent)
	if err != nil {
		log.Printf("UpdateFileContentError: %s", err.Error())
		return errors.New("UpdateFileContentInfraError")
	}

	log.Printf(
		"File '%s' content updated.",
		sourcePath.GetFileName().String(),
	)

	return nil
}

func (uc UpdateUnixFiles) Execute(
	updateUnixFiles dto.UpdateUnixFiles,
) (dto.UpdateProcessReport, error) {
	updateProcessReport := dto.NewUpdateProcessReport(
		[]valueObject.UnixFilePath{},
		[]valueObject.UpdateProcessFailure{},
	)

	for _, sourcePath := range updateUnixFiles.SourcePaths {
		if updateUnixFiles.Permissions != nil {
			err := uc.updateFilePermissions(sourcePath, *updateUnixFiles.Permissions)
			if err != nil {
				updateFailure, err := uc.updateFailureFactory(sourcePath, err.Error())
				if err != nil {
					log.Printf("AddUpdatePermissionsFailureError: %s", err.Error())
				}

				updateProcessReport.FailedPathsWithReason = append(
					updateProcessReport.FailedPathsWithReason,
					updateFailure,
				)
				continue
			}
		}

		if updateUnixFiles.DestinationPath != nil {
			err := uc.moveFile(sourcePath, *updateUnixFiles.DestinationPath)
			if err != nil {
				updateFailure, err := uc.updateFailureFactory(sourcePath, err.Error())
				if err != nil {
					log.Printf("AddMoveFailureError: %s", err.Error())
				}

				updateProcessReport.FailedPathsWithReason = append(
					updateProcessReport.FailedPathsWithReason,
					updateFailure,
				)
				continue
			}
		}

		if updateUnixFiles.EncodedContent != nil {
			err := uc.updateFileContent(sourcePath, *updateUnixFiles.EncodedContent)
			if err != nil {
				updateFailure, err := uc.updateFailureFactory(sourcePath, err.Error())
				if err != nil {
					log.Printf("AddUpdateContentFailureError: %s", err.Error())
				}

				updateProcessReport.FailedPathsWithReason = append(
					updateProcessReport.FailedPathsWithReason,
					updateFailure,
				)
				continue
			}
		}

		updateProcessReport.FilePathsSuccessfullyUpdated = append(
			updateProcessReport.FilePathsSuccessfullyUpdated,
			sourcePath,
		)
	}

	return updateProcessReport, nil
}
