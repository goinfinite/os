package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func CompressUnixFiles(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	compressUnixFiles dto.CompressUnixFiles,
) (dto.CompressionProcessReport, error) {
	fileDestinationPath := compressUnixFiles.DestinationPath

	compressionProcessReport, err := filesCmdRepo.Compress(compressUnixFiles)
	if err != nil {
		return compressionProcessReport, err
	}

	allPathsFailedCompression := len(compressionProcessReport.Failure) == len(compressUnixFiles.Paths)
	if allPathsFailedCompression {
		log.Printf(
			"UnableToCompressFilesAndDirectories: File compressed %s wasn't created.",
			fileDestinationPath,
		)
		return compressionProcessReport, errors.New("UnableToCompressFilesAndDirectories")
	}

	log.Printf("File compressed %s created.", fileDestinationPath)

	return compressionProcessReport, nil
}
