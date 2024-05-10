package valueObject

import (
	"errors"
	"slices"
	"strings"
)

type DataFieldType string

var validDataFieldTypes = []string{
	"checkbox",
	"color",
	"date",
	"email",
	"image",
	"number",
	"password",
	"radio",
	"range",
	"search",
	"tel",
	"text",
	"time",
	"url",
}

func NewDataFieldType(value string) (DataFieldType, error) {
	value = strings.ToLower(value)
	if !slices.Contains(validDataFieldTypes, value) {
		return "", errors.New("InvalidDataFieldType")
	}

	return DataFieldType(value), nil
}

func NewDataFieldTypePanic(value string) DataFieldType {
	vo, err := NewDataFieldType(value)
	if err != nil {
		panic(err)
	}

	return vo
}

func (vo DataFieldType) String() string {
	return string(vo)
}
