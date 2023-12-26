package valueObject

type UploadProcessFailure struct {
	FilePath UnixFilePath `json:"filePath"`
	Reason   string       `json:"reason"`
}

func NewUploadProcessFailure(
	filePath UnixFilePath,
	reason string,
) UploadProcessFailure {
	return UploadProcessFailure{
		FilePath: filePath,
		Reason:   reason,
	}
}
