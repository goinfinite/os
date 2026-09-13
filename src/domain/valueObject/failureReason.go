package valueObject

import (
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

const malformedFailureReason FailureReason = "MalformedFailureReason"

const maxFailureReasonStringLength = 2048

type FailureReason string

func NewFailureReason(value any) FailureReason {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return malformedFailureReason
	}

	if len(stringValue) == 0 {
		return malformedFailureReason
	}

	if len(stringValue) > maxFailureReasonStringLength {
		stringValue = stringValue[:maxFailureReasonStringLength]
	}

	return FailureReason(stringValue)
}

func (vo FailureReason) String() string {
	return string(vo)
}
