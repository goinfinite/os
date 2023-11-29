package valueObject

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

const unixFilePathRegexExpression = `^\/?[^\n\r\t\f\0\?\[\]\<\>]+$`
const unixFileRelativePathRegexExpression = `(\.\.\/)|^\.\/|^\/\.\/`

type UnixFilePath string

func NewUnixFilePath(value string) (UnixFilePath, error) {
	isFileName := len(value) > 0 && !strings.Contains(value, "/")
	if isFileName {
		value = "/" + value
	}

	unixFilePath := UnixFilePath(value)

	if !unixFilePath.isValid() {
		return "", errors.New("InvalidUnixFilePath")
	}

	if unixFilePath.isRelative() {
		return "", errors.New("RelativePathNotAllowed")
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

func (unixFilePath UnixFilePath) isRelative() bool {
	unixFileRelativePathRegex := regexp.MustCompile(unixFileRelativePathRegexExpression)
	return unixFileRelativePathRegex.MatchString(string(unixFilePath))
}

func (unixFilePath UnixFilePath) GetFileName() (UnixFileName, error) {
	var unixFileName UnixFileName

	isDir := strings.HasSuffix(string(unixFilePath), "/")
	if isDir {
		return unixFileName, errors.New("UnableToGetFileName")
	}

	unixFileBase := filepath.Base(string(unixFilePath))
	unixFileName, err := NewUnixFileName(unixFileBase)

	return unixFileName, err
}

func (unixFilePath UnixFilePath) GetFileExtension() (UnixFileExtension, error) {
	return NewUnixFileExtension(filepath.Ext(string(unixFilePath)))
}

func (unixFilePath UnixFilePath) GetFileDir() (UnixFilePath, error) {
	return NewUnixFilePath(filepath.Dir(string(unixFilePath)))
}

func (unixFilePath UnixFilePath) String() string {
	return string(unixFilePath)
}
