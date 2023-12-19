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
		log.Printf("PathIsDirError: %s", err)
		return uploadProcessInfo, errors.New("PathIsDirError")
	}

	if !fileIsDir {
		return uploadProcessInfo, errors.New("PathIsFile")
	}

	for _, multipartFile := range uploadUnixFiles.MultipartFiles {
		multipartFileSizeInGB := multipartFile.GetFileSize().ToGiB()
		if multipartFileSizeInGB > 5 {
			errMessage := "File size in greater than 5 GB"
			log.Printf("UploadFileError: %s", errMessage)

			uploadProcessFailure := uploadProcessFailureFactory(
				multipartFile.GetFileName(),
				errMessage,
			)
			uploadProcessInfo.Failure = append(
				uploadProcessInfo.Failure,
				uploadProcessFailure,
			)

			continue
		}

		err := filesCmdRepo.Upload(fileDestinationPath, multipartFile)
		if err != nil {
			log.Printf("UploadFileError: %v", err)

			uploadProcessFailure := uploadProcessFailureFactory(
				multipartFile.GetFileName(),
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
			multipartFile.GetFileName().String(),
		)

		log.Printf(
			"File '%s' content upload to '%s'.",
			multipartFile.GetFileName().String(),
			fileDestinationPath.String(),
		)
	}

	return uploadProcessInfo, nil
}
