package valueObject

import (
	"errors"
)

type UnixFileContent string

const FileContentMaxSizeInMb = 5

func NewUnixFileContent(value string) (UnixFileContent, error) {
	charsAmountIn1Mb := 1048576
	contentLimitInCharsAmount := charsAmountIn1Mb * FileContentMaxSizeInMb
	if len(value) > contentLimitInCharsAmount {
		return "", errors.New("InvalidUnixFileContent")
	}

	return UnixFileContent(value), nil
}

func NewUnixFileContentPanic(value string) UnixFileContent {
	unixFileContent, err := NewUnixFileContent(value)
	if err != nil {
		panic(err)
	}

	return unixFileContent
}

func (unixFileContent UnixFileContent) String() string {
	return string(unixFileContent)
}
