package valueObject

type UpdateProcessFailure struct {
	FilePath UnixFilePath `json:"filePath"`
	Reason   string       `json:"reason"`
}

func NewUpdateProcessFailure(
	filePath UnixFilePath,
	reason string,
) UpdateProcessFailure {
	return UpdateProcessFailure{
		FilePath: filePath,
		Reason:   reason,
	}
}
