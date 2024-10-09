package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type UpdateProcessReport struct {
	FilePathsSuccessfullyUpdated []valueObject.UnixFilePath         `json:"filePathsSuccessfullyUpdated"`
	FailedPathsWithReason        []valueObject.UpdateProcessFailure `json:"failedPathsWithReason"`
}

func NewUpdateProcessReport(
	filePathsSuccessfullyUpdated []valueObject.UnixFilePath,
	failedPathsWithReason []valueObject.UpdateProcessFailure,
) UpdateProcessReport {
	return UpdateProcessReport{
		FilePathsSuccessfullyUpdated: filePathsSuccessfullyUpdated,
		FailedPathsWithReason:        failedPathsWithReason,
	}
}
