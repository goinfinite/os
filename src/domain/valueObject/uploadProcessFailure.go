package valueObject

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

type UploadProcessFailure struct {
	FileName tkValueObject.UnixFileName `json:"fileName"`
	Reason   FailureReason              `json:"reason"`
}

func NewUploadProcessFailure(
	fileName tkValueObject.UnixFileName,
	reason FailureReason,
) UploadProcessFailure {
	return UploadProcessFailure{
		FileName: fileName,
		Reason:   reason,
	}
}
