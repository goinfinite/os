package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func UploadUnixFiles(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	uploadDto dto.UploadUnixFiles,
) (dto.UploadProcessReport, error) {
	maxFileSizeInGb := int64(5)

	tooBigFiles := []valueObject.UploadProcessFailure{}
	filesToUpload := []valueObject.FileStreamHandler{}

	for _, fileStream := range uploadDto.FileStreamHandlers {
		fileSizeInGb := fileStream.Size.ToGiB()
		if fileSizeInGb < maxFileSizeInGb {
			filesToUpload = append(filesToUpload, fileStream)
			continue
		}

		failureReason, _ := valueObject.NewFailureReason("FileTooBig")
		processFailure := valueObject.NewUploadProcessFailure(
			fileStream.Name,
			failureReason,
		)
		tooBigFiles = append(tooBigFiles, processFailure)

		log.Printf("FileTooBig: %s", fileStream.Name)
	}

	uploadDto.FileStreamHandlers = filesToUpload

	uploadProcessReport, err := filesCmdRepo.Upload(uploadDto)
	if err != nil {
		log.Printf("UploadUnixFileInfraError: %s", err.Error())
		return uploadProcessReport, errors.New("UploadUnixFileInfraError")
	}

	uploadProcessReport.FailedNamesWithReason = append(
		uploadProcessReport.FailedNamesWithReason,
		tooBigFiles...,
	)

	log.Printf("Files uploaded to '%s'.", uploadDto.DestinationPath)

	return uploadProcessReport, nil
}
