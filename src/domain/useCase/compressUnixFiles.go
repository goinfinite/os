package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
)

func CompressUnixFiles(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo,
	compressDto dto.CompressUnixFiles,
) (dto.CompressionProcessReport, error) {
	compressionProcessReport, err := filesCmdRepo.Compress(compressDto)
	if err != nil {
		slog.Error("CompressUnixFilesError", slog.String("err", err.Error()))
		return compressionProcessReport, errors.New("CompressUnixFilesInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).CompressUnixFile(compressDto)

	NormalizeKnownUnixFilePathPermissions(filesCmdRepo, compressDto.DestinationPath)

	return compressionProcessReport, nil
}
