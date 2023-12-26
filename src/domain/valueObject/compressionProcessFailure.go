package valueObject

type CompressionProcessFailure struct {
	FilePath UnixFilePath       `json:"filePath"`
	Reason   ProcessFileFailure `json:"reason"`
}

func NewCompressionProcessFailure(
	filePath UnixFilePath,
	reason ProcessFileFailure,
) CompressionProcessFailure {
	return CompressionProcessFailure{
		FilePath: filePath,
		Reason:   reason,
	}
}
