package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type UnixCompressionType string

var validUnixCompressionTypes = []string{
	"gzip",
	"zip",
}

func NewUnixCompressionType(value string) (UnixCompressionType, error) {
	value = strings.ToLower(value)
	if !slices.Contains(validUnixCompressionTypes, value) {
		return "", errors.New("InvalidUnixCompressionType")
	}
	return UnixCompressionType(value), nil
}

func NewUnixCompressionTypePanic(value string) UnixCompressionType {
	unixCompressionType, err := NewUnixCompressionType(value)
	if err != nil {
		panic(err)
	}
	return unixCompressionType
}

func (unixCompressionType UnixCompressionType) String() string {
	return string(unixCompressionType)
}
