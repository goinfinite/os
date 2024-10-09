package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func CompressUnixFiles(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	compressUnixFiles dto.CompressUnixFiles,
) (dto.CompressionProcessReport, error) {
	compressionProcessReport, err := filesCmdRepo.Compress(compressUnixFiles)
	if err != nil {
		log.Printf("CompressUnixFilesInfraError: %s", err.Error())
		return compressionProcessReport, errors.New("CompressUnixFilesInfraError")
	}

	log.Printf("Compressed file '%s' created.", compressUnixFiles.DestinationPath)

	return compressionProcessReport, nil
}
