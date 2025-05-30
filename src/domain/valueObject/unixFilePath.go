package valueObject

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const unixFilePathRegexExpression = `^\/?[^\n\r\t\f\0\?\[\]\<\>]+$`
const unixFileRelativePathRegexExpression = `\.\.\/|^\.\/|^\/\.\/`

type UnixFilePath string

const UnixFilePathFileSystemRootDir = UnixFilePath("/")
const UnixFilePathAppWorkingDir = UnixFilePath("/app")
const UnixFilePathTrashDir = UnixFilePath("/app/.trash")

func NewUnixFilePath(value interface{}) (filePath UnixFilePath, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return filePath, errors.New("UnixFilePathValueMustBeString")
	}

	unixFilePathRegex := regexp.MustCompile(unixFilePathRegexExpression)
	if !unixFilePathRegex.MatchString(stringValue) {
		return filePath, errors.New("InvalidUnixFilePath")
	}

	isOnlyFileName := !strings.Contains(stringValue, "/")
	if isOnlyFileName {
		return filePath, errors.New("PathIsFileNameOnly")
	}

	unixFileRelativePathRegex := regexp.MustCompile(unixFileRelativePathRegexExpression)
	if unixFileRelativePathRegex.MatchString(stringValue) {
		return filePath, errors.New("RelativePathNotAllowed")
	}

	if stringValue != UnixFilePathFileSystemRootDir.String() && strings.HasSuffix(stringValue, "/") {
		stringValue = strings.TrimSuffix(stringValue, "/")
	}

	return UnixFilePath(stringValue), nil
}

func (vo UnixFilePath) ReadWithoutExtension() UnixFilePath {
	unixFilePathExtStr := filepath.Ext(string(vo))
	if unixFilePathExtStr == "" {
		return vo
	}

	unixFilePathWithoutExtStr := strings.TrimSuffix(string(vo), unixFilePathExtStr)
	unixFilePathWithoutExt, _ := NewUnixFilePath(unixFilePathWithoutExtStr)
	return unixFilePathWithoutExt
}

func (vo UnixFilePath) ReadFileName() UnixFileName {
	if vo.IsFileSystemRootDir() {
		return UnixFileName("/")
	}

	unixFileBase := filepath.Base(string(vo))
	unixFileName, _ := NewUnixFileName(unixFileBase)
	return unixFileName
}

func (vo UnixFilePath) ReadFileNameWithoutExtension() UnixFileName {
	unixFileBase := filepath.Base(string(vo))
	unixFilePathExt := filepath.Ext(string(vo))
	unixFileBaseWithoutExtStr := strings.TrimSuffix(string(unixFileBase), unixFilePathExt)
	unixFileNameWithoutExt, _ := NewUnixFileName(unixFileBaseWithoutExtStr)
	return unixFileNameWithoutExt
}

func (vo UnixFilePath) ReadFileExtension() (UnixFileExtension, error) {
	unixFileExtensionStr := filepath.Ext(string(vo))
	return NewUnixFileExtension(unixFileExtensionStr)
}

func (vo UnixFilePath) ReadFileDir() UnixFilePath {
	unixFileDirPath, _ := NewUnixFilePath(filepath.Dir(string(vo)))
	return unixFileDirPath
}

func (vo UnixFilePath) String() string {
	return string(vo)
}

func (vo UnixFilePath) IsFileSystemRootDir() bool {
	return vo == UnixFilePathFileSystemRootDir
}
