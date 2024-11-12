package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type UpdateUnixFiles struct {
	filesCmdRepo          repository.FilesCmdRepo
	activityRecordCmdRepo repository.ActivityRecordCmdRepo
}

func NewUpdateUnixFiles(
	filesCmdRepo repository.FilesCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
) UpdateUnixFiles {
	return UpdateUnixFiles{
		filesCmdRepo:          filesCmdRepo,
		activityRecordCmdRepo: activityRecordCmdRepo,
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

	return valueObject.NewUpdateProcessFailure(filePath, failureReason), nil
}

func (uc UpdateUnixFiles) updateFilePermissions(
	sourcePath valueObject.UnixFilePath,
	permissions valueObject.UnixFilePermissions,
) error {
	err := uc.filesCmdRepo.UpdatePermissions(sourcePath, permissions)
	if err != nil {
		slog.Error("UpdateFilePermissionsError", slog.Any("err", err))
		return errors.New("UpdateFilePermissionsInfraError")
	}

	return nil
}

func (uc UpdateUnixFiles) moveFile(
	sourcePath valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
) error {
	shouldOverwrite := false
	err := uc.filesCmdRepo.Move(
		sourcePath, destinationPath, shouldOverwrite,
	)
	if err != nil {
		slog.Error("MoveFileError", slog.Any("err", err))
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
		slog.Error("UpdateFileContentError", slog.Any("err", err))
		return errors.New("UpdateFileContentInfraError")
	}

	return nil
}

func (uc UpdateUnixFiles) Execute(
	updateDto dto.UpdateUnixFiles,
) (dto.UpdateProcessReport, error) {
	updateProcessReport := dto.NewUpdateProcessReport(
		[]valueObject.UnixFilePath{},
		[]valueObject.UpdateProcessFailure{},
	)

	for _, sourcePath := range updateDto.SourcePaths {
		if updateDto.Permissions != nil {
			err := uc.updateFilePermissions(sourcePath, *updateDto.Permissions)
			if err != nil {
				updateFailure, err := uc.updateFailureFactory(sourcePath, err.Error())
				if err != nil {
					slog.Error("AddUpdatePermissionsFailureError", slog.Any("err", err))
				}

				updateProcessReport.FailedPathsWithReason = append(
					updateProcessReport.FailedPathsWithReason,
					updateFailure,
				)
				continue
			}
		}

		if updateDto.EncodedContent != nil {
			err := uc.updateFileContent(sourcePath, *updateDto.EncodedContent)
			if err != nil {
				updateFailure, err := uc.updateFailureFactory(sourcePath, err.Error())
				if err != nil {
					slog.Error("AddUpdateContentFailureError", slog.Any("err", err))
				}

				updateProcessReport.FailedPathsWithReason = append(
					updateProcessReport.FailedPathsWithReason,
					updateFailure,
				)
				continue
			}
		}

		if updateDto.DestinationPath != nil {
			err := uc.moveFile(sourcePath, *updateDto.DestinationPath)
			if err != nil {
				updateFailure, err := uc.updateFailureFactory(sourcePath, err.Error())
				if err != nil {
					slog.Error("AddMoveFailureError", slog.Any("err", err))
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

	NewCreateSecurityActivityRecord(uc.activityRecordCmdRepo).
		UpdateUnixFiles(updateDto)

	return updateProcessReport, nil
}
