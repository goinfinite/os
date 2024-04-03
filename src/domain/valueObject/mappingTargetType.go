package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type MappingTargetType string

var ValidMappingTargetTypes = []string{
	"url",
	"service",
	"response-code",
	"inline-html",
	"static-files",
}

func NewMappingTargetType(value string) (MappingTargetType, error) {
	value = strings.ToLower(value)
	if !slices.Contains(ValidMappingTargetTypes, value) {
		return "", errors.New("InvalidMappingTargetType")
	}
	return MappingTargetType(value), nil
}

func NewMappingTargetTypePanic(value string) MappingTargetType {
	mtt, err := NewMappingTargetType(value)
	if err != nil {
		panic(err)
	}
	return mtt
}

func (mtt MappingTargetType) String() string {
	return string(mtt)
}

func (mttPtr *MappingTargetType) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	unquotedValue := strings.Trim(valueStr, "\"")

	mtt, err := NewMappingTargetType(unquotedValue)
	if err != nil {
		return err
	}

	*mttPtr = mtt
	return nil
}
