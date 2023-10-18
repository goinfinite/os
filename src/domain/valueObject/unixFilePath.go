package valueObject

import (
	"errors"
	"regexp"
)

type UnixFilePath string

func NewUnixFilePath(unixFilePathStr string) (UnixFilePath, error) {
	unixFilePath := UnixFilePath(unixFilePathStr)
	if !unixFilePath.isValid() {
		return "", errors.New("InvalidUnixFilePath")
	}
	return unixFilePath, nil
}

func (unixFilePath UnixFilePath) isValid() bool {
	unixFilePathRegexExpression := "^\\S*\\/[^\\/]*\\.[A-z]{1,5}$"
	unixFilePathRegexRegex := regexp.MustCompile(unixFilePathRegexExpression)
	return unixFilePathRegexRegex.MatchString(string(unixFilePath))
}

func (unixFilePath UnixFilePath) String() string {
	return string(unixFilePath)
}
