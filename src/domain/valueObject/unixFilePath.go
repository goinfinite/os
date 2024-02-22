package valueObject

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

const unixFilePathRegexExpression = `^\/?[^\n\r\t\f\0\?\[\]\<\>]+$`
const unixFileRelativePathRegexExpression = `\.\.\/|^\.\/|^\/\.\/`

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

func (unixFilePath UnixFilePath) GetWithoutExtension() UnixFilePath {
	unixFilePathExtStr := filepath.Ext(string(unixFilePath))
	unixFilePathWithoutExtStr := strings.TrimSuffix(string(unixFilePath), unixFilePathExtStr)
	unixFilePathWithoutExt, _ := NewUnixFilePath(unixFilePathWithoutExtStr)
	return unixFilePathWithoutExt
}

func (unixFilePath UnixFilePath) GetFileName() UnixFileName {
	unixFileBase := filepath.Base(string(unixFilePath))
	unixFileName, _ := NewUnixFileName(unixFileBase)
	return unixFileName
}

func (unixFilePath UnixFilePath) GetFileNameWithoutExtension() UnixFileName {
	unixFileBase := filepath.Base(string(unixFilePath))
	unixFilePathExt := filepath.Ext(string(unixFilePath))
	unixFileBaseWithoutExtStr := strings.TrimSuffix(string(unixFileBase), unixFilePathExt)
	unixFileNameWithoutExt, _ := NewUnixFileName(unixFileBaseWithoutExtStr)
	return unixFileNameWithoutExt
}

func (unixFilePath UnixFilePath) GetFileExtension() (UnixFileExtension, error) {
	unixFileExtensionStr := filepath.Ext(string(unixFilePath))
	return NewUnixFileExtension(unixFileExtensionStr)
}

func (unixFilePath UnixFilePath) GetFileDir() UnixFilePath {
	unixFileDirPath, _ := NewUnixFilePath(filepath.Dir(string(unixFilePath)))
	return unixFileDirPath
}

func (unixFilePath UnixFilePath) String() string {
	return string(unixFilePath)
}
