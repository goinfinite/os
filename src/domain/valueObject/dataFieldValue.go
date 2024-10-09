package valueObject

import (
	"errors"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type DataFieldValue string

func NewDataFieldValue(value interface{}) (
	dataFieldValue DataFieldValue, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return dataFieldValue, errors.New("DataFieldValueMustBeString")
	}

	if len(stringValue) == 0 {
		return dataFieldValue, errors.New("EmptyDataFieldValue")
	}

	if len(stringValue) >= 2048 {
		return dataFieldValue, errors.New("DataFieldValueTooBig")
	}

	return DataFieldValue(stringValue), nil
}

func (vo DataFieldValue) String() string {
	return string(vo)
}
