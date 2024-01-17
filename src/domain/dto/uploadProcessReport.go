package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UploadProcessReport struct {
	FileNamesSuccessfullyUploaded []valueObject.UnixFileName         `json:"fileNamesSuccessfullyUploaded"`
	FailedNamesWithReason         []valueObject.UploadProcessFailure `json:"failedNamesWithReason"`
	DestinationPath               valueObject.UnixFilePath           `json:"destinationPath"`
}

func NewUploadProcessReport(
	fileNamesSuccessfullyUploaded []valueObject.UnixFileName,
	failedNamesWithReason []valueObject.UploadProcessFailure,
	destinationPath valueObject.UnixFilePath,
) UploadProcessReport {
	return UploadProcessReport{
		FileNamesSuccessfullyUploaded: fileNamesSuccessfullyUploaded,
		FailedNamesWithReason:         failedNamesWithReason,
		DestinationPath:               destinationPath,
	}
}
