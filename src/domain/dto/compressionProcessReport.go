package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CompressionProcessReport struct {
	FilePathsSuccessfullyCompressed []tkValueObject.UnixAbsoluteFilePath    `json:"filePathsSuccessfullyCompressed"`
	FailedPathsWithReason           []valueObject.CompressionProcessFailure `json:"failedPathsWithReason"`
	DestinationPath                 tkValueObject.UnixAbsoluteFilePath      `json:"destinationPath"`
}

func NewCompressionProcessReport(
	filePathsSuccessfullyCompressed []tkValueObject.UnixAbsoluteFilePath,
	failedPathsWithReason []valueObject.CompressionProcessFailure,
	destinationPath tkValueObject.UnixAbsoluteFilePath,
) CompressionProcessReport {
	return CompressionProcessReport{
		FilePathsSuccessfullyCompressed: filePathsSuccessfullyCompressed,
		FailedPathsWithReason:           failedPathsWithReason,
		DestinationPath:                 destinationPath,
	}
}
