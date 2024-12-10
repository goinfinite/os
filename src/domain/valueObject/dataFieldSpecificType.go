package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type DataFieldSpecificType string

var validDataFieldSpecificTypes = []string{
	"password", "username", "email",
}

func NewDataFieldSpecificType(value interface{}) (
	dataFieldSpecificType DataFieldSpecificType, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return dataFieldSpecificType, errors.New("DataFieldSpecificTypeMustBeString")
	}

	stringValue = strings.ToLower(stringValue)
	if !slices.Contains(validDataFieldSpecificTypes, stringValue) {
		return dataFieldSpecificType, errors.New("InvalidDataFieldSpecificType")
	}

	return DataFieldSpecificType(stringValue), nil
}

func (vo DataFieldSpecificType) String() string {
	return string(vo)
}
