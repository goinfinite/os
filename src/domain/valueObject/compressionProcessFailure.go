package valueObject

type CompressionProcessFailure struct {
	FilePath UnixFilePath `json:"filePath"`
	Reason   string       `json:"reason"`
}

func NewCompressionProcessFailure(
	filePath UnixFilePath,
	reason string,
) CompressionProcessFailure {
	return CompressionProcessFailure{
		FilePath: filePath,
		Reason:   reason,
	}
}
