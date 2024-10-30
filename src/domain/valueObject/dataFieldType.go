package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type DataFieldType string

var validDataFieldTypes = []string{
	"checkbox", "color", "date", "email", "image", "number", "password", "radio",
	"range", "search", "select", "tel", "text", "time", "url",
}

func NewDataFieldType(value interface{}) (
	dataFieldType DataFieldType, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return dataFieldType, errors.New("DataFieldTypeMustBeString")
	}

	stringValue = strings.ToLower(stringValue)
	if !slices.Contains(validDataFieldTypes, stringValue) {
		return dataFieldType, errors.New("InvalidDataFieldType")
	}

	return DataFieldType(stringValue), nil
}

func (vo DataFieldType) String() string {
	return string(vo)
}
