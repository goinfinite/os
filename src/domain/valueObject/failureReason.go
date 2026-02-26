package valueObject

import (
	"errors"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type FailureReason string

func NewFailureReason(value interface{}) (failureReason FailureReason, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return failureReason, errors.New("FailureReasonMustBeString")
	}

	if len(stringValue) == 0 {
		return failureReason, errors.New("EmptyFailureReason")
	}

	if len(stringValue) > 2048 {
		stringValue = stringValue[:2048]
	}

	return FailureReason(stringValue), nil
}

func (vo FailureReason) String() string {
	return string(vo)
}
