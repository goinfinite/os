package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UploadProcessReport struct {
	FilePathsSuccessfullyUploaded         []valueObject.UnixFileName         `json:"filePathsSuccessfullyUploaded"`
	FilePathsThatFailedToUploadWithReason []valueObject.UploadProcessFailure `json:"filePathsThatFailedToUploadWithReason"`
	DestinationPath                       valueObject.UnixFilePath           `json:"destinationPath"`
}

func NewUploadProcessReport(
	filePathsSuccessfullyUploaded []valueObject.UnixFileName,
	filePathsThatFailedToUploadWithReason []valueObject.UploadProcessFailure,
	destinationPath valueObject.UnixFilePath,
) UploadProcessReport {
	return UploadProcessReport{
		FilePathsSuccessfullyUploaded:         filePathsSuccessfullyUploaded,
		FilePathsThatFailedToUploadWithReason: filePathsThatFailedToUploadWithReason,
		DestinationPath:                       destinationPath,
	}
}
