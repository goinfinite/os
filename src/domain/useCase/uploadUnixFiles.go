package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

type UploadProcessFailure struct {
	FileName string `json:"fileName"`
	Reason   string `json:"reason"`
}

type UploadProcessInfo struct {
	Success     []string               `json:"success"`
	Failure     []UploadProcessFailure `json:"failure"`
	Destination string                 `json:"destination"`
}

func uploadProcessFailureFactory(
	fileName valueObject.UnixFileName,
	errMessage string,
) UploadProcessFailure {
	return UploadProcessFailure{
		FileName: fileName.String(),
		Reason:   errMessage,
	}
}

func UploadUnixFiles(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	uploadUnixFiles dto.UploadUnixFiles,
) (UploadProcessInfo, error) {
	fileDestinationPath := uploadUnixFiles.DestinationPath

	uploadProcessInfo := UploadProcessInfo{
		Success:     []string{},
		Failure:     []UploadProcessFailure{},
		Destination: fileDestinationPath.String(),
	}

	fileIsDir, err := filesQueryRepo.IsDir(fileDestinationPath)
	if err != nil {
		return uploadProcessInfo, err
	}

	if !fileIsDir {
		return uploadProcessInfo, errors.New("PathCannotBeAFile")
	}

	for _, fileStreamHandler := range uploadUnixFiles.FileStreamHandlers {
		fileStreamHandlerSizeInGB := fileStreamHandler.GetFileSize().ToGiB()
		if fileStreamHandlerSizeInGB > 5 {
			errMessage := "File size in greater than 5 GB"
			log.Printf("UploadFileError: %s", errMessage)

			uploadProcessFailure := uploadProcessFailureFactory(
				fileStreamHandler.GetFileName(),
				errMessage,
			)
			uploadProcessInfo.Failure = append(
				uploadProcessInfo.Failure,
				uploadProcessFailure,
			)

			continue
		}

		err := filesCmdRepo.Upload(fileDestinationPath, fileStreamHandler)
		if err != nil {
			uploadProcessFailure := uploadProcessFailureFactory(
				fileStreamHandler.GetFileName(),
				err.Error(),
			)
			uploadProcessInfo.Failure = append(
				uploadProcessInfo.Failure,
				uploadProcessFailure,
			)

			continue
		}

		uploadProcessInfo.Success = append(
			uploadProcessInfo.Success,
			fileStreamHandler.GetFileName().String(),
		)

		log.Printf(
			"File '%s' content upload to '%s'.",
			fileStreamHandler.GetFileName().String(),
			fileDestinationPath.String(),
		)
	}

	return uploadProcessInfo, nil
}
