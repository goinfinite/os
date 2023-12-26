package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CompressionProcessReport struct {
	FilePathsSuccessfullyCompressed []valueObject.UnixFilePath              `json:"filePathsSuccessfullyCompressed"`
	FailedPathsWithReason           []valueObject.CompressionProcessFailure `json:"failedPathsWithReason"`
	DestinationPath                 valueObject.UnixFilePath                `json:"destinationPath"`
}

func NewCompressionProcessReport(
	filePathsSuccessfullyCompressed []valueObject.UnixFilePath,
	failedPathsWithReason []valueObject.CompressionProcessFailure,
	destinationPath valueObject.UnixFilePath,
) CompressionProcessReport {
	return CompressionProcessReport{
		FilePathsSuccessfullyCompressed: filePathsSuccessfullyCompressed,
		FailedPathsWithReason:           failedPathsWithReason,
		DestinationPath:                 destinationPath,
	}
}
