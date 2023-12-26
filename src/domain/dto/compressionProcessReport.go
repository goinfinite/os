package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CompressionProcessReport struct {
	FilePathsSuccessfullyCompressed         []valueObject.UnixFilePath              `json:"filePathsSuccessfullyCompressed"`
	FilePathsThatFailedToCompressWithReason []valueObject.CompressionProcessFailure `json:"filePathsThatFailedToCompressWithReason"`
	DestinationPath                         valueObject.UnixFilePath                `json:"destinationPath"`
}

func NewCompressionProcessReport(
	filePathsSuccessfullyCompressed []valueObject.UnixFilePath,
	filePathsThatFailedToCompressWithReason []valueObject.CompressionProcessFailure,
	destinationPath valueObject.UnixFilePath,
) CompressionProcessReport {
	return CompressionProcessReport{
		FilePathsSuccessfullyCompressed:         filePathsSuccessfullyCompressed,
		FilePathsThatFailedToCompressWithReason: filePathsThatFailedToCompressWithReason,
		DestinationPath:                         destinationPath,
	}
}
