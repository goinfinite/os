package valueObject

type UploadProcessFailure struct {
	FileName UnixFileName       `json:"fileName"`
	Reason   ProcessFileFailure `json:"reason"`
}

func NewUploadProcessFailure(
	fileName UnixFileName,
	reason ProcessFileFailure,
) UploadProcessFailure {
	return UploadProcessFailure{
		FileName: fileName,
		Reason:   reason,
	}
}
