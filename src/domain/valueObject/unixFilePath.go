package valueObject

import (
	"errors"
	"regexp"
)

const unixFilePathRegexExpression = `^\/(?:[\w\p{Latin}\. \-]+\/)*[\w\p{Latin}\. \-]+$`

type UnixFilePath string

func NewUnixFilePath(value string) (UnixFilePath, error) {
	unixFilePath := UnixFilePath(value)
	if !unixFilePath.isValid() {
		return "", errors.New("InvalidUnixFilePath")
	}
	return unixFilePath, nil
}

func NewUnixFilePathPanic(value string) UnixFilePath {
	unixFilePath, err := NewUnixFilePath(value)
	if err != nil {
		panic(err)
	}
	return unixFilePath
}

func (unixFilePath UnixFilePath) isValid() bool {
	unixFilePathRegex := regexp.MustCompile(unixFilePathRegexExpression)
	return unixFilePathRegex.MatchString(string(unixFilePath))
}

func (unixFilePath UnixFilePath) String() string {
	return string(unixFilePath)
}
