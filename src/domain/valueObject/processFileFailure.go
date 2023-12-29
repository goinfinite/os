package valueObject

import (
	"errors"
)

type FileProcessingFailure string

func NewFileProcessingFailure(value string) (FileProcessingFailure, error) {
	fileProcessingFailure := FileProcessingFailure(value)
	if !fileProcessingFailure.isValid() {
		return "", errors.New("InvalidFileProcessingFailure")
	}
	return fileProcessingFailure, nil
}

func NewFileProcessingFailurePanic(value string) FileProcessingFailure {
	fileProcessingFailure, err := NewFileProcessingFailure(value)
	if err != nil {
		panic(err)
	}
	return fileProcessingFailure
}

func (fileProcessingFailure FileProcessingFailure) isValid() bool {
	isTooShort := len(string(fileProcessingFailure)) < 1

	size5MBInBytes := (1024 * 1024) * 5
	isTooLong := len(string(fileProcessingFailure)) > size5MBInBytes

	return !isTooShort && !isTooLong
}

func (fileProcessingFailure FileProcessingFailure) String() string {
	return string(fileProcessingFailure)
}
