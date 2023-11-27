package valueObject

import (
	"errors"
	"regexp"
	"slices"
)

const unixFileNameRegexExpression = `^[^\n\r\t\f\0\?\[\]\<\>\/]{1,256}$`

var reservedUnixFileNames = []string{".", "..", "*", "/", "\\"}

type UnixFileName string

func NewUnixFileName(unixFileNameStr string) (UnixFileName, error) {
	unixFileName := UnixFileName(unixFileNameStr)
	if !unixFileName.isValid() {
		return "", errors.New("InvalidUnixFileName")
	}
	return unixFileName, nil
}

func NewUnixFileNamePanic(unixFileNameStr string) UnixFileName {
	unixFileName, err := NewUnixFileName(unixFileNameStr)
	if err != nil {
		panic(err)
	}
	return unixFileName
}

func (unixFileName UnixFileName) isValid() bool {
	unixFileNameRegexRegex := regexp.MustCompile(unixFileNameRegexExpression)
	isValidFormat := unixFileNameRegexRegex.MatchString(string(unixFileName))

	isReservedUnixFileName := slices.Contains(reservedUnixFileNames, string(unixFileName))

	return isValidFormat && !isReservedUnixFileName
}

func (unixFileName UnixFileName) String() string {
	return string(unixFileName)
}
