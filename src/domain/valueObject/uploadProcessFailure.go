package valueObject

type UploadProcessFailure struct {
	FileName UnixFileName `json:"fileName"`
	Reason   string       `json:"reason"`
}

func NewUploadProcessFailure(
	fileName UnixFileName,
	reason string,
) UploadProcessFailure {
	return UploadProcessFailure{
		FileName: fileName,
		Reason:   reason,
	}
}
