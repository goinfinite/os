package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
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

	_, err := filesQueryRepo.Get(fileDestinationPath)
	if err != nil {

	}

	fileIsDir, err := fileDestinationPath.IsDir()
	if err != nil {
		log.Printf("PathIsDirError: %s", err)
		return uploadProcessInfo, errors.New("PathIsDirError")
	}

	if !fileIsDir {
		return uploadProcessInfo, errors.New("PathIsFile")
	}

	for _, multipartFile := range uploadUnixFiles.MultipartFiles {
		err := filesCmdRepo.Upload(fileDestinationPath, multipartFile)
		if err != nil {
			log.Printf("UploadFileError: %v", err)

			uploadProcessFailure := UploadProcessFailure{
				FileName: multipartFile.GetFileName().String(),
				Reason:   err.Error(),
			}
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
