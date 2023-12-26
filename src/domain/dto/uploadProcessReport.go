package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UploadProcessReport struct {
	Success     []valueObject.UnixFilePath              `json:"success"`
	Failure     []valueObject.CompressionProcessFailure `json:"failure"`
	Destination valueObject.UnixFilePath                `json:"destination"`
}

func NewUploadProcessReport(
	success []valueObject.UnixFilePath,
	failure []valueObject.CompressionProcessFailure,
	destination valueObject.UnixFilePath,
) UploadProcessReport {
	return UploadProcessReport{
		Success:     success,
		Failure:     failure,
		Destination: destination,
	}
}
