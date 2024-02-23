package valueObject

import (
	"errors"
)

type UnixFileContent string

const charsAmountIn1Mb = 1048576
const fileContentLimitInMb = 5
const FileContentMaxSize = charsAmountIn1Mb * fileContentLimitInMb

func NewUnixFileContent(value string) (UnixFileContent, error) {
	isTooShort := len(value) < 1
	isTooBig := len(value) > FileContentMaxSize

	if isTooShort || isTooBig {
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
