package valueObject

import "errors"

type FileProcessingFailure string

func NewFileProcessingFailure(value string) (FileProcessingFailure, error) {
	maxProcessingFailureSize := 256

	if len(value) < 1 {
		return "", errors.New("EmptyProcessingFailure")
	}

	if len(value) > maxProcessingFailureSize {
		maxProcessingFailureSizeIndex := maxProcessingFailureSize - 1
		partialProcessingFailure := value[:maxProcessingFailureSizeIndex]
		value = partialProcessingFailure
	}

	return FileProcessingFailure(value), nil
}

func NewFileProcessingFailurePanic(value string) FileProcessingFailure {
	fileProcessingFailure, err := NewFileProcessingFailure(value)
	if err != nil {
		panic(err)
	}
	return fileProcessingFailure
}

func (fileProcessingFailure FileProcessingFailure) String() string {
	return string(fileProcessingFailure)
}
