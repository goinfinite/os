package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
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
) valueObject.UpdateProcessFailure {
	failureReason, err := valueObject.NewFailureReason(errMessage)
	if err != nil {
		slog.Debug(err.Error(), slog.String("errMessage", errMessage))
		failureReason, _ = valueObject.NewFailureReason("MalformedFailureReason")
	}

	return valueObject.NewUpdateProcessFailure(filePath, failureReason)
}

func (uc UpdateUnixFiles) updateFilePermissions(
	sourcePath valueObject.UnixFilePath,
	permissions valueObject.UnixFilePermissions,
) error {
	updatePermissions := dto.NewUpdateUnixFilePermissions(
		sourcePath, permissions, nil,
	)

	err := uc.filesCmdRepo.UpdatePermissions(updatePermissions)
	if err != nil {
		slog.Error("UpdateFilePermissionsError", slog.String("err", err.Error()))
		return errors.New("UpdateFilePermissionsInfraError")
	}

	return nil
}

func (uc UpdateUnixFiles) updateFileContent(
	sourcePath valueObject.UnixFilePath,
	encodedContent valueObject.EncodedContent,
) error {
	updateContentDto := dto.NewUpdateUnixFileContent(sourcePath, encodedContent)

	err := uc.filesCmdRepo.UpdateContent(updateContentDto)
	if err != nil {
		slog.Error("UpdateFileContentError", slog.String("err", err.Error()))
		return errors.New("UpdateFileContentInfraError")
	}

	return nil
}

func (uc UpdateUnixFiles) updateFileOwnership(
	sourcePath valueObject.UnixFilePath,
	ownership valueObject.UnixFileOwnership,
) error {
	updateOwnershipDto := dto.NewUpdateUnixFileOwnership(sourcePath, ownership)

	err := uc.filesCmdRepo.UpdateOwnership(updateOwnershipDto)
	if err != nil {
		slog.Error("UpdateFileOwnershipError", slog.String("err", err.Error()))
		return errors.New("UpdateFileOwnershipInfraError")
	}

	return nil
}

func (uc UpdateUnixFiles) fixFilePermissions(
	sourcePath valueObject.UnixFilePath,
) error {
	sourcePathStr := sourcePath.String()
	if infraHelper.FileExists(sourcePathStr) {
		return errors.New("FileOrDirectoryNotFound")
	}

	filePermissions := valueObject.NewUnixFileDefaultPermissions()

	dirPermissions := valueObject.NewUnixDirDefaultPermissions()
	if sourcePathStr == "/app/html" {
		dirPermissions, _ = valueObject.NewUnixFilePermissions("777")
	}

	updatePermissionsDto := dto.NewUpdateUnixFilePermissions(
		sourcePath, filePermissions, &dirPermissions,
	)

	err := uc.filesCmdRepo.UpdatePermissions(updatePermissionsDto)
	if err != nil {
		slog.Error("FixFilePermissionsError", slog.String("err", err.Error()))
		return errors.New("FixFilePermissionsInfraError")
	}

	if sourcePathStr == "/app" {
		defaultOwnership := valueObject.NewUnixFileDefaultOwnership()
		updateOwnershipDto := dto.NewUpdateUnixFileOwnership(
			sourcePath, defaultOwnership,
		)

		err = uc.filesCmdRepo.UpdateOwnership(updateOwnershipDto)
		if err != nil {
			return errors.New("FixAppDirOwnershipError: " + err.Error())
		}
	}

	return nil
}

func (uc UpdateUnixFiles) moveFile(
	sourcePath valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
) error {
	shouldOverwrite := false
	moveDto := dto.NewMoveUnixFile(sourcePath, destinationPath, shouldOverwrite)

	err := uc.filesCmdRepo.Move(moveDto)
	if err != nil {
		slog.Error("MoveFileError", slog.String("err", err.Error()))
		return errors.New("MoveFileInfraError")
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
				updateFailure := uc.updateFailureFactory(sourcePath, err.Error())
				updateProcessReport.FailedPathsWithReason = append(
					updateProcessReport.FailedPathsWithReason, updateFailure,
				)
				continue
			}
		}

		if updateDto.Ownership != nil {
			err := uc.updateFileOwnership(sourcePath, *updateDto.Ownership)
			if err != nil {
				updateFailure := uc.updateFailureFactory(sourcePath, err.Error())
				updateProcessReport.FailedPathsWithReason = append(
					updateProcessReport.FailedPathsWithReason, updateFailure,
				)
				continue
			}
		}

		if updateDto.EncodedContent != nil {
			err := uc.updateFileContent(sourcePath, *updateDto.EncodedContent)
			if err != nil {
				updateFailure := uc.updateFailureFactory(sourcePath, err.Error())
				updateProcessReport.FailedPathsWithReason = append(
					updateProcessReport.FailedPathsWithReason, updateFailure,
				)
				continue
			}
		}

		if updateDto.ShouldFixPermissions != nil {
			err := uc.fixFilePermissions(sourcePath)
			if err != nil {
				updateFailure := uc.updateFailureFactory(sourcePath, err.Error())
				updateProcessReport.FailedPathsWithReason = append(
					updateProcessReport.FailedPathsWithReason, updateFailure,
				)
				continue
			}
		}

		if updateDto.DestinationPath != nil {
			err := uc.moveFile(sourcePath, *updateDto.DestinationPath)
			if err != nil {
				updateFailure := uc.updateFailureFactory(sourcePath, err.Error())
				updateProcessReport.FailedPathsWithReason = append(
					updateProcessReport.FailedPathsWithReason, updateFailure,
				)
				continue
			}
		}

		updateProcessReport.FilePathsSuccessfullyUpdated = append(
			updateProcessReport.FilePathsSuccessfullyUpdated, sourcePath,
		)
	}

	NewCreateSecurityActivityRecord(uc.activityRecordCmdRepo).
		UpdateUnixFiles(updateDto)

	return updateProcessReport, nil
}
