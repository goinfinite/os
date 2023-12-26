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
	uploadUnixFiles dto.UploadUnixFiles,
) (dto.UploadProcessReport, error) {
	filesLargerThanAllowed := []valueObject.FileStreamHandler{}
	filesWithAllowedSizes := []valueObject.FileStreamHandler{}
	largerFileErrMessage := "File size is greater than 5 GB"
	for _, fileToUploadStream := range uploadUnixFiles.FileStreamHandlers {
		fileStreamHandlerSizeInGB := fileToUploadStream.GetFileSize().ToGiB()
		if fileStreamHandlerSizeInGB > 5 {
			log.Printf("UploadUnixFileError: %s", largerFileErrMessage)

			filesLargerThanAllowed = append(filesLargerThanAllowed, fileToUploadStream)

			continue
		}

		filesWithAllowedSizes = append(filesWithAllowedSizes, fileToUploadStream)
	}

	uploadUnixFiles.FileStreamHandlers = filesWithAllowedSizes

	uploadProcessReport := filesCmdRepo.Upload(uploadUnixFiles)

	for _, largeFile := range filesLargerThanAllowed {
		uploadProcessReport.FilePathsThatFailedToUploadWithReason = append(
			uploadProcessReport.FilePathsThatFailedToUploadWithReason,
			valueObject.NewUploadProcessFailure(largeFile.GetFileName(), largerFileErrMessage),
		)
	}

	allFilesFailedToUpload := len(uploadProcessReport.FilePathsThatFailedToUploadWithReason) == len(uploadUnixFiles.FileStreamHandlers)
	if allFilesFailedToUpload {
		return uploadProcessReport, errors.New("UploadUnixFileInfraError")
	}

	log.Printf("Files uploaded to '%s'.", uploadUnixFiles.DestinationPath)

	return uploadProcessReport, nil
}
