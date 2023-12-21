package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CompressionProcessReport struct {
	Success     []valueObject.UnixFilePath              `json:"success"`
	Failure     []valueObject.CompressionProcessFailure `json:"failure"`
	Destination valueObject.UnixFilePath                `json:"destination"`
}

func NewCompressionProcessReport(
	success []valueObject.UnixFilePath,
	failure []valueObject.CompressionProcessFailure,
	destination valueObject.UnixFilePath,
) CompressionProcessReport {
	return CompressionProcessReport{
		Success:     success,
		Failure:     failure,
		Destination: destination,
	}
}
