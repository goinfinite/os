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
	copyUnixFile dto.CopyUnixFile,
) error {
	err := filesCmdRepo.Copy(copyUnixFile)
	if err != nil {
		slog.Error("CopyUnixFileInfraError", slog.Any("err", err))
		return errors.New("CopyUnixFileInfraError")
	}

	return nil
}
