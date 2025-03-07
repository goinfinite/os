package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

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

	return nil
}
