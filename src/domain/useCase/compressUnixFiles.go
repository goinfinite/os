package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

type CompressionProcessFailure struct {
	FilePath string `json:"filePath"`
	Reason   string `json:"reason"`
}

type CompressionProcessInfo struct {
	Success     []string                    `json:"success"`
	Failure     []CompressionProcessFailure `json:"failure"`
	Destination string                      `json:"destination"`
}

func compressionProcessFailureFactory(
	filePath valueObject.UnixFilePath,
	err error,
) CompressionProcessFailure {
	return CompressionProcessFailure{
		FilePath: filePath.String(),
		Reason:   err.Error(),
	}
}

func CompressUnixFiles(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	compressUnixFiles dto.CompressUnixFiles,
) (CompressionProcessInfo, error) {
	fileDestinationPath := compressUnixFiles.DestinationPath

	compressionProcessInfo := CompressionProcessInfo{
		Success:     []string{},
		Failure:     []CompressionProcessFailure{},
		Destination: fileDestinationPath.String(),
	}

	unixFiles, _ := filesQueryRepo.Get(fileDestinationPath)

	if len(unixFiles) > 0 {
		log.Print("PathAlreadyExists")

		return compressionProcessInfo, errors.New("PathAlreadyExists")
	}

	var filesToCompress []valueObject.UnixFilePath
	for _, filePath := range compressUnixFiles.Paths {
		unixDestinationFiles, err := filesQueryRepo.Get(filePath)

		if err != nil || len(unixDestinationFiles) < 1 {
			log.Printf("PathDoesNotExists: %v", err)

			compressionProcessFailure := compressionProcessFailureFactory(filePath, err)
			compressionProcessInfo.Failure = append(
				compressionProcessInfo.Failure,
				compressionProcessFailure,
			)

			continue
		}

		filesToCompress = append(filesToCompress, filePath)
	}

	allPathsFailedInCompression := len(compressionProcessInfo.Failure) == len(compressUnixFiles.Paths)
	if allPathsFailedInCompression {
		log.Printf(
			"UnableToCompressFilesAndDirectories: File compressed %s wasn't created.",
			fileDestinationPath,
		)
		return compressionProcessInfo, errors.New("UnableToCompressFilesAndDirectories")
	}

	err := filesCmdRepo.Compress(
		filesToCompress,
		fileDestinationPath,
		compressUnixFiles.CompressionType,
	)
	if err != nil {
		log.Printf("UnableToCompressFilesAndDirectories: %s", err.Error())

		for _, filePath := range filesToCompress {
			compressionProcessFailure := compressionProcessFailureFactory(filePath, err)
			compressionProcessInfo.Failure = append(
				compressionProcessInfo.Failure,
				compressionProcessFailure,
			)
		}

		return compressionProcessInfo, errors.New("UnableToCompressFilesAndDirectories")
	}

	for _, filePath := range filesToCompress {
		compressionProcessInfo.Success = append(
			compressionProcessInfo.Success,
			filePath.String(),
		)
	}

	log.Printf("File compressed %s created.", fileDestinationPath)

	return compressionProcessInfo, nil
}
