package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UpdateProcessReport struct {
	FilePathsSuccessfullyUpdated []tkValueObject.UnixAbsoluteFilePath `json:"filePathsSuccessfullyUpdated"`
	FailedPathsWithReason        []valueObject.UpdateProcessFailure   `json:"failedPathsWithReason"`
}

func NewUpdateProcessReport(
	filePathsSuccessfullyUpdated []tkValueObject.UnixAbsoluteFilePath,
	failedPathsWithReason []valueObject.UpdateProcessFailure,
) UpdateProcessReport {
	return UpdateProcessReport{
		FilePathsSuccessfullyUpdated: filePathsSuccessfullyUpdated,
		FailedPathsWithReason:        failedPathsWithReason,
	}
}
