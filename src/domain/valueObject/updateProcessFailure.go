package valueObject

type UpdateProcessFailure struct {
	FilePath UnixFilePath          `json:"filePath"`
	Reason   FileProcessingFailure `json:"reason"`
}

func NewUpdateProcessFailure(
	filePath UnixFilePath,
	reason FileProcessingFailure,
) UpdateProcessFailure {
	return UpdateProcessFailure{
		FilePath: filePath,
		Reason:   reason,
	}
}
