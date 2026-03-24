package valueObject

import (
	"errors"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type UnixFileContent string

const FileContentMaxSizeInMb = 5

func NewUnixFileContent(value interface{}) (
	unixFileContent UnixFileContent, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return unixFileContent, errors.New("UnixFileContentMustBeString")
	}

	charsAmountIn1Mb := 1048576
	contentLimitInCharsAmount := charsAmountIn1Mb * FileContentMaxSizeInMb
	if len(stringValue) > contentLimitInCharsAmount {
		return unixFileContent, errors.New("UnixFileContentTooBig")
	}

	return UnixFileContent(stringValue), nil
}

func (vo UnixFileContent) String() string {
	return string(vo)
}
