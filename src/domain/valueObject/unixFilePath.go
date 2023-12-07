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
	unixFilePathStr := string(unixFilePath)

	isOnlyFileName := !strings.Contains(unixFilePathStr, "/")

	unixFileRelativePathRegex := regexp.MustCompile(unixFileRelativePathRegexExpression)
	return isOnlyFileName || unixFileRelativePathRegex.MatchString(unixFilePathStr)
}

func (unixFilePath UnixFilePath) GetFileName() (UnixFileName, error) {
	var unixFileName UnixFileName

	isDir := strings.HasSuffix(string(unixFilePath), "/")
	if isDir {
		return unixFileName, errors.New("FilePathIsDir")
	}

	unixFileBase := filepath.Base(string(unixFilePath))
	unixFileName, err := NewUnixFileName(unixFileBase)

	return unixFileName, err
}

func (unixFilePath UnixFilePath) GetFileExtension() (UnixFileExtension, error) {
	unixFileExtensionStr := filepath.Ext(string(unixFilePath))
	if len(unixFileExtensionStr) < 1 {
		return "", nil
	}

	return NewUnixFileExtension(unixFileExtensionStr)
}

func (unixFilePath UnixFilePath) GetFileDir() (UnixFilePath, error) {
	return NewUnixFilePath(filepath.Dir(string(unixFilePath)))
}

func (unixFilePath UnixFilePath) String() string {
	return string(unixFilePath)
}
