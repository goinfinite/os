package valueObject

import (
	"errors"
	"regexp"
	"slices"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const unixFileNameRegexExpression = `^[^\n\r\t\f\0\?\[\]\<\>\/]{1,512}$`

var reservedUnixFileNames = []string{".", "..", "*", "/", "\\"}

type UnixFileName string

func NewUnixFileName(value interface{}) (fileName UnixFileName, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return fileName, errors.New("UnixFileNameValueMustBeString")
	}

	re := regexp.MustCompile(unixFileNameRegexExpression)
	if !re.MatchString(stringValue) {
		return fileName, errors.New("InvalidUnixFileName")
	}

	if slices.Contains(reservedUnixFileNames, stringValue) {
		return fileName, errors.New("ReservedUnixFileName")
	}

	return UnixFileName(stringValue), nil
}

func (vo UnixFileName) String() string {
	return string(vo)
}
