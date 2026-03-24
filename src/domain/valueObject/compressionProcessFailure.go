package valueObject

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

type CompressionProcessFailure struct {
	FilePath tkValueObject.UnixAbsoluteFilePath `json:"filePath"`
	Reason   FailureReason                      `json:"reason"`
}

func NewCompressionProcessFailure(
	filePath tkValueObject.UnixAbsoluteFilePath,
	reason FailureReason,
) CompressionProcessFailure {
	return CompressionProcessFailure{
		FilePath: filePath,
		Reason:   reason,
	}
}
