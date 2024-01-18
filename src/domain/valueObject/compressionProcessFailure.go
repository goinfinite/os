package valueObject

type CompressionProcessFailure struct {
	FilePath UnixFilePath          `json:"filePath"`
	Reason   FileProcessingFailure `json:"reason"`
}

func NewCompressionProcessFailure(
	filePath UnixFilePath,
	reason FileProcessingFailure,
) CompressionProcessFailure {
	return CompressionProcessFailure{
		FilePath: filePath,
		Reason:   reason,
	}
}
