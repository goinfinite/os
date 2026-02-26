package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UploadProcessReport struct {
	FileNamesSuccessfullyUploaded []tkValueObject.UnixFileName       `json:"fileNamesSuccessfullyUploaded"`
	FailedNamesWithReason         []valueObject.UploadProcessFailure `json:"failedNamesWithReason"`
	DestinationPath               tkValueObject.UnixAbsoluteFilePath `json:"destinationPath"`
}

func NewUploadProcessReport(
	fileNamesSuccessfullyUploaded []tkValueObject.UnixFileName,
	failedNamesWithReason []valueObject.UploadProcessFailure,
	destinationPath tkValueObject.UnixAbsoluteFilePath,
) UploadProcessReport {
	return UploadProcessReport{
		FileNamesSuccessfullyUploaded: fileNamesSuccessfullyUploaded,
		FailedNamesWithReason:         failedNamesWithReason,
		DestinationPath:               destinationPath,
	}
}
