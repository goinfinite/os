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
	errMessage string,
) CompressionProcessFailure {
	return CompressionProcessFailure{
		FilePath: filePath.String(),
		Reason:   errMessage,
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

	unixFileExists, err := filesQueryRepo.Exists(fileDestinationPath)
	if err != nil {
		return compressionProcessInfo, err
	}

	if unixFileExists {
		return compressionProcessInfo, errors.New("PathAlreadyExists")
	}

	var filesToCompress []valueObject.UnixFilePath
	for _, filePath := range compressUnixFiles.Paths {
		unixDestinationFileExists, err := filesQueryRepo.Exists(filePath)
		if !unixDestinationFileExists {
			errMessage := "PathDoesNotExists"
			if err != nil {
				errMessage = err.Error()
			}

			compressionProcessFailure := compressionProcessFailureFactory(filePath, errMessage)
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

	err = filesCmdRepo.Compress(
		filesToCompress,
		fileDestinationPath,
		compressUnixFiles.CompressionType,
	)
	if err != nil {
		log.Printf("UnableToCompressFilesAndDirectories: %s", err.Error())

		for _, filePath := range filesToCompress {
			compressionProcessFailure := compressionProcessFailureFactory(filePath, err.Error())
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
