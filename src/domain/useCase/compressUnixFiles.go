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
) error {
	compressionErrorCount := 0

	for _, filePath := range compressUnixFiles.Paths {
		unixFiles, _ := filesQueryRepo.Get(compressUnixFiles.DestinationPath)

		if len(unixFiles) > 0 {
			compressionErrorCount++

			log.Print("PathAlreadyExists")
			continue
		}

		unixDestinationFiles, err := filesQueryRepo.Get(filePath)

		if err != nil || len(unixDestinationFiles) < 1 {
			compressionErrorCount++

			log.Printf("PathDoesNotExists: %v", err)
			continue
		}
	}

	err := filesCmdRepo.Compress(
		compressUnixFiles.Paths,
		compressUnixFiles.DestinationPath,
		compressUnixFiles.CompressionType,
	)
	if err != nil {
		log.Printf("CompressError: %s", err.Error())
		return errors.New("UnableToCompressFilesAndDirectories")
	}

	if compressionErrorCount == len(compressUnixFiles.Paths) {
		log.Printf(
			"UnableToCompressFilesAndDirectories: File compressed %s wasn't created.",
			compressUnixFiles.DestinationPath,
		)
		return errors.New("UnableToCompressFilesAndDirectories")
	}

	log.Printf("File compressed %s created.", compressUnixFiles.DestinationPath)

	return nil
}
