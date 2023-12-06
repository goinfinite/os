package valueObject

import (
	"errors"
)

type UnixFileType string

func NewUnixFileType(value string) (UnixFileType, error) {
	unixFileType := UnixFileType(value)
	if !unixFileType.isValid() {
		return "", errors.New("InvalidUnixFileType")
	}

	return unixFileType, nil
}

func NewUnixFileTypePanic(value string) UnixFileType {
	unixFileType, err := NewUnixFileType(value)
	if err != nil {
		panic(err)
	}
	return unixFileType
}

func (unixFileType UnixFileType) isValid() bool {
	return unixFileType == "directory" || unixFileType == "file"
}

func (unixFileType UnixFileType) String() string {
	return string(unixFileType)
}
