package valueObject

type UploadProcessFailure struct {
	FileName UnixFileName  `json:"fileName"`
	Reason   FailureReason `json:"reason"`
}

func NewUploadProcessFailure(
	fileName UnixFileName,
	reason FailureReason,
) UploadProcessFailure {
	return UploadProcessFailure{
		FileName: fileName,
		Reason:   reason,
	}
}
