package valueObject

type UpdateProcessFailure struct {
	FilePath UnixFilePath  `json:"filePath"`
	Reason   FailureReason `json:"reason"`
}

func NewUpdateProcessFailure(
	filePath UnixFilePath,
	reason FailureReason,
) UpdateProcessFailure {
	return UpdateProcessFailure{
		FilePath: filePath,
		Reason:   reason,
	}
}
