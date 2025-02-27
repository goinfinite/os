package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func UploadUnixFiles(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	uploadDto dto.UploadUnixFiles,
) (dto.UploadProcessReport, error) {
	maxFileSizeInGb := int64(5)

	tooBigFiles := []valueObject.UploadProcessFailure{}
	filesToUpload := []valueObject.FileStreamHandler{}

	failureReasonStr := "FileTooBig"
	for _, fileStream := range uploadDto.FileStreamHandlers {
		fileSizeInGb := fileStream.Size.ToGiB()
		if fileSizeInGb < maxFileSizeInGb {
			filesToUpload = append(filesToUpload, fileStream)
			continue
		}

		failureReason, _ := valueObject.NewFailureReason(failureReasonStr)
		processFailure := valueObject.NewUploadProcessFailure(
			fileStream.Name, failureReason,
		)
		tooBigFiles = append(tooBigFiles, processFailure)

		slog.Debug(failureReasonStr, slog.String("fileName", fileStream.Name.String()))
	}

	uploadDto.FileStreamHandlers = filesToUpload

	uploadProcessReport, err := filesCmdRepo.Upload(uploadDto)
	if err != nil {
		slog.Error("UploadUnixFileError", slog.String("err", err.Error()))
		return uploadProcessReport, errors.New("UploadUnixFileInfraError")
	}

	uploadProcessReport.FailedNamesWithReason = append(
		uploadProcessReport.FailedNamesWithReason, tooBigFiles...,
	)

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).UploadUnixFiles(uploadDto)

	return uploadProcessReport, nil
}
