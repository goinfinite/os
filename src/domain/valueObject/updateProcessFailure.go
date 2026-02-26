package valueObject

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

type UpdateProcessFailure struct {
	FilePath tkValueObject.UnixAbsoluteFilePath `json:"filePath"`
	Reason   FailureReason                      `json:"reason"`
}

func NewUpdateProcessFailure(
	filePath tkValueObject.UnixAbsoluteFilePath,
	reason FailureReason,
) UpdateProcessFailure {
	return UpdateProcessFailure{
		FilePath: filePath,
		Reason:   reason,
	}
}
