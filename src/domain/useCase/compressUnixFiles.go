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
		unixFiles, _ := filesQueryRepo.Get(filePath)

		if len(unixFiles) > 0 {
			return errors.New("PathAlreadyExists")
		}

		unixDestinationFiles, err := filesQueryRepo.Get(compressUnixFiles.DestinationPath)

		if err != nil || len(unixDestinationFiles) < 1 {
			log.Printf("PathDoesNotExists: %v", err)
			continue
		}

		isDir, err := filePath.IsDir()
		if err != nil {
			log.Printf("PathIsDirError: %s", err)
			continue
		}

		inodeName := "File"
		if isDir {
			inodeName = "Directory"
		}

		err = filesCmdRepo.Compress(
			filePath,
			compressUnixFiles.DestinationPath,
			compressUnixFiles.CompressionType,
		)
		if err != nil {
			compressionErrorCount++

			log.Printf("%sCompressError: %s", inodeName, err.Error())
			continue
		}
	}

	if compressionErrorCount == len(compressUnixFiles.Paths) {
		log.Printf(
			"UnableToCompressFilesAndDirectories: File compressed %s  wasn't created.",
			compressUnixFiles.DestinationPath,
		)

		return errors.New("UnableToCompressFilesAndDirectories")
	}

	log.Printf("File compressed %s created.", compressUnixFiles.DestinationPath)

	return nil
}
