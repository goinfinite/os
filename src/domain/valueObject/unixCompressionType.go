package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type UnixCompressionType string

var ValidUnixCompressionTypes = []string{
	"tgz", "zip",
}

func NewUnixCompressionType(value interface{}) (
	unixCompressionType UnixCompressionType, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return unixCompressionType, errors.New("UnixCompressionTypeMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidUnixCompressionTypes, stringValue) {
		return unixCompressionType, errors.New("InvalidUnixCompressionType")
	}

	return UnixCompressionType(stringValue), nil
}

func (vo UnixCompressionType) String() string {
	return string(vo)
}
