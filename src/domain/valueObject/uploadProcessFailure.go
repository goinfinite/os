package valueObject

type UploadProcessFailure struct {
	FileName UnixFileName          `json:"fileName"`
	Reason   FileProcessingFailure `json:"reason"`
}

func NewUploadProcessFailure(
	fileName UnixFileName,
	reason FileProcessingFailure,
) UploadProcessFailure {
	return UploadProcessFailure{
		FileName: fileName,
		Reason:   reason,
	}
}
