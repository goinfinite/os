package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func CopyUnixFile(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	copyDto dto.CopyUnixFile,
) error {
	err := filesCmdRepo.Copy(copyDto)
	if err != nil {
		slog.Error("CopyUnixFileInfraError", slog.String("err", err.Error()))
		return errors.New("CopyUnixFileInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).CopyUnixFile(copyDto)

	return nil
}
