package valueObject

type CompressionProcessFailure struct {
	FilePath UnixFilePath  `json:"filePath"`
	Reason   FailureReason `json:"reason"`
}

func NewCompressionProcessFailure(
	filePath UnixFilePath,
	reason FailureReason,
) CompressionProcessFailure {
	return CompressionProcessFailure{
		FilePath: filePath,
		Reason:   reason,
	}
}
