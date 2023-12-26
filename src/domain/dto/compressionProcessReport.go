package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CompressionProcessReport struct {
	FilePathsSuccessfullyProcessed         []valueObject.UnixFilePath              `json:"filePathsSuccessfullyProcessed"`
	FilePathsThatFailedToProcessWithReason []valueObject.CompressionProcessFailure `json:"filePathsThatFailedToProcessWithReason"`
	DestinationPath                        valueObject.UnixFilePath                `json:"destinationPath"`
}

func NewCompressionProcessReport(
	filePathsSuccessfullyProcessed []valueObject.UnixFilePath,
	filePathsThatFailedToProcessWithReason []valueObject.CompressionProcessFailure,
	destinationPath valueObject.UnixFilePath,
) CompressionProcessReport {
	return CompressionProcessReport{
		FilePathsSuccessfullyProcessed:         filePathsSuccessfullyProcessed,
		FilePathsThatFailedToProcessWithReason: filePathsThatFailedToProcessWithReason,
		DestinationPath:                        destinationPath,
	}
}
