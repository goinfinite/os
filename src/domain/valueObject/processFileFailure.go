package valueObject

import (
	"errors"
)

type ProcessFileFailure string

func NewProcessFileFailure(value string) (ProcessFileFailure, error) {
	processFileFailure := ProcessFileFailure(value)
	if !processFileFailure.isValid() {
		return "", errors.New("InvalidProcessFileFailure")
	}
	return processFileFailure, nil
}

func NewProcessFileFailurePanic(value string) ProcessFileFailure {
	processFileFailure, err := NewProcessFileFailure(value)
	if err != nil {
		panic(err)
	}
	return processFileFailure
}

func (processFileFailure ProcessFileFailure) isValid() bool {
	isTooShort := len(string(processFileFailure)) < 1

	size5MBInBytes := (1024 * 1024) * 5
	isTooLong := len(string(processFileFailure)) > size5MBInBytes

	return !isTooShort && !isTooLong
}

func (processFileFailure ProcessFileFailure) String() string {
	return string(processFileFailure)
}
