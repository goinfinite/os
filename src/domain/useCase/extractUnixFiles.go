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
	extractDto dto.ExtractUnixFiles,
) error {
	err := filesCmdRepo.Extract(extractDto)
	if err != nil {
		slog.Error("ExtractUnixFilesError", slog.Any("err", err))
		return errors.New("ExtractUnixFilesInfraError")
	}

	return nil
}
