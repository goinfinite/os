package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UploadProcessReport struct {
	Success     []valueObject.UnixFileName         `json:"success"`
	Failure     []valueObject.UploadProcessFailure `json:"failure"`
	Destination valueObject.UnixFilePath           `json:"destination"`
}

func NewUploadProcessReport(
	success []valueObject.UnixFileName,
	failure []valueObject.UploadProcessFailure,
	destination valueObject.UnixFilePath,
) UploadProcessReport {
	return UploadProcessReport{
		Success:     success,
		Failure:     failure,
		Destination: destination,
	}
}
