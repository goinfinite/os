package valueObject

import (
	"errors"
)

type UnixFileContent string

func NewUnixFileContent(value string) (UnixFileContent, error) {
	unixFileContent := UnixFileContent(value)
	if !unixFileContent.isValid() {
		return "", errors.New("InvalidUnixFileContent")
	}

	return unixFileContent, nil
}

func NewUnixFileContentPanic(value string) UnixFileContent {
	unixFileContent, err := NewUnixFileContent(value)
	if err != nil {
		panic(err)
	}

	return unixFileContent
}

func (unixFileContent UnixFileContent) isValid() bool {
	unixFileContentStr := string(unixFileContent)

	isTooShort := len(unixFileContentStr) < 1

	contentCharAmountIn1Mb := 1048576
	contentLimitInMb := 5
	contentLimitInCharAmount := contentCharAmountIn1Mb * contentLimitInMb
	isTooBig := len(unixFileContentStr) > contentLimitInCharAmount

	return !isTooShort && !isTooBig
}

func (unixFileContent UnixFileContent) String() string {
	return string(unixFileContent)
}
