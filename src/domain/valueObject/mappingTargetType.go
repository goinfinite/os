package valueObject

import (
	"errors"

	"golang.org/x/exp/slices"
)

type MappingTargetType string

var ValidMappingTargetTypes = []string{
	"url",
	"service",
	"responseCode",
}

func NewMappingTargetType(value string) (MappingTargetType, error) {
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
