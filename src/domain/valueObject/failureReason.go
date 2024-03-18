package valueObject

import "errors"

type FailureReason string

func NewFailureReason(value string) (FailureReason, error) {
	maxProcessingFailureSize := 256

	if len(value) < 1 {
		return "", errors.New("EmptyFailureReason")
	}

	if len(value) > maxProcessingFailureSize {
		maxProcessingFailureSizeIndex := maxProcessingFailureSize - 1
		partialProcessingFailure := value[:maxProcessingFailureSizeIndex]
		value = partialProcessingFailure
	}

	return FailureReason(value), nil
}

func NewFailureReasonPanic(value string) FailureReason {
	fileProcessingFailure, err := NewFailureReason(value)
	if err != nil {
		panic(err)
	}
	return fileProcessingFailure
}

func (fileProcessingFailure FailureReason) String() string {
	return string(fileProcessingFailure)
}
