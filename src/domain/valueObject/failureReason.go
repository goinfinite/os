package valueObject

import (
	"errors"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type FailureReason string

func NewFailureReason(value interface{}) (failureReason FailureReason, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return failureReason, errors.New("FailureReasonMustBeString")
	}

	if len(stringValue) == 0 {
		return failureReason, errors.New("FailureReasonEmpty")
	}

	maxProcessingFailureSize := 256
	if len(stringValue) > 256 {
		maxProcessingFailureSizeIndex := maxProcessingFailureSize - 1
		partialProcessingFailure := stringValue[:maxProcessingFailureSizeIndex]
		stringValue = partialProcessingFailure
	}

	return FailureReason(stringValue), nil
}

func (vo FailureReason) String() string {
	return string(vo)
}
