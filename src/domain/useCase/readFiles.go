package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func ReadFiles(
	filesQueryRepo repository.FilesQueryRepo,
	unixFilePath valueObject.UnixFilePath,
) ([]entity.UnixFile, error) {
	filesList, err := filesQueryRepo.Read(unixFilePath)
	if err != nil {
		slog.Error("ReadFilesError", slog.String("err", err.Error()))
		return filesList, errors.New("ReadFilesInfraError")
	}

	return filesList, nil
}
