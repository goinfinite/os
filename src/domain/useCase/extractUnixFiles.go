package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func ExtractUnixFiles(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	extractDto dto.ExtractUnixFiles,
) error {
	err := filesCmdRepo.Extract(extractDto)
	if err != nil {
		slog.Error("ExtractUnixFilesInfraError", slog.Any("err", err))
		return errors.New("ExtractUnixFilesInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).ExtractUnixFile(extractDto)

	return nil
}
