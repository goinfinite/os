package valueObject

import (
	"errors"
	"regexp"
	"slices"
	"strings"

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

	stringValue = strings.TrimSpace(stringValue)

	unixFileNameRegex := regexp.MustCompile(unixFileNameRegexExpression)
	if !unixFileNameRegex.MatchString(stringValue) {
		return "", errors.New("InvalidUnixFileName")
	}

	if slices.Contains(reservedUnixFileNames, stringValue) {
		return "", errors.New("ReservedUnixFileName")
	}

	return UnixFileName(stringValue), nil
}

func (vo UnixFileName) String() string {
	return string(vo)
}
