package valueObject

type FileProcessingFailure string

func NewFileProcessingFailure(value string) (FileProcessingFailure, error) {
	maxProcessingFailureSize := 256

	maxProcessingFailureSizeIndex := maxProcessingFailureSize - 1
	partialProcessingFailure := value[:maxProcessingFailureSizeIndex]
	return FileProcessingFailure(partialProcessingFailure), nil
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
