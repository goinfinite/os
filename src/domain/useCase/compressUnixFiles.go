package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func CompressUnixFiles(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	compressDto dto.CompressUnixFiles,
) (dto.CompressionProcessReport, error) {
	compressionProcessReport, err := filesCmdRepo.Compress(compressDto)
	if err != nil {
		slog.Error("CompressUnixFilesError", slog.Any("err", err))
		return compressionProcessReport, errors.New("CompressUnixFilesInfraError")
	}

	return compressionProcessReport, nil
}
