package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UploadProcessReport struct {
	FilePathsSuccessfullyUploaded []valueObject.UnixFileName         `json:"filePathsSuccessfullyUploaded"`
	FailedPathsWithReason         []valueObject.UploadProcessFailure `json:"failedPathsWithReason"`
	DestinationPath               valueObject.UnixFilePath           `json:"destinationPath"`
}

func NewUploadProcessReport(
	filePathsSuccessfullyUploaded []valueObject.UnixFileName,
	failedPathsWithReason []valueObject.UploadProcessFailure,
	destinationPath valueObject.UnixFilePath,
) UploadProcessReport {
	return UploadProcessReport{
		FilePathsSuccessfullyUploaded: filePathsSuccessfullyUploaded,
		FailedPathsWithReason:         failedPathsWithReason,
		DestinationPath:               destinationPath,
	}
}
