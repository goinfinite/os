package useCase

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func NormalizeKnownUnixFilePathPermissions(
	filesCmdRepo repository.FilesCmdRepo,
	filePath valueObject.UnixFilePath,
) {
	isAppDirectoryDescendant := strings.HasPrefix(
		filePath.String(),
		valueObject.UnixFilePathAppWorkingDir.String(),
	)
	if !isAppDirectoryDescendant {
		return
	}

	err := filesCmdRepo.UpdateOwnership(dto.NewUpdateUnixFileOwnership(
		filePath, valueObject.UnixFileOwnershipAppWorkingDir,
	))
	if err != nil {
		slog.Debug(
			"UpdateOwnershipInfraError",
			slog.String("method", "NormalizeKnownUnixFilePathPermissions"),
			slog.String("err", err.Error()),
		)
	}
}

func CreateUnixFile(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateUnixFile,
) error {
	err := filesCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateUnixFileInfraError", slog.String("err", err.Error()))
		return errors.New("CreateUnixFileInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateUnixFile(createDto)

	NormalizeKnownUnixFilePathPermissions(filesCmdRepo, createDto.FilePath)

	return nil
}
