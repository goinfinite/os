package valueObject

import (
	"errors"
	"slices"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type UnixCompressionType string

var ValidUnixCompressionTypes = []string{
	"tgz", "zip",
}

func NewUnixCompressionType(value interface{}) (
	unixCompressionType UnixCompressionType, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
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
