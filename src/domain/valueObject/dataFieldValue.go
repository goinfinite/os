package valueObject

import (
	"errors"
	"strings"

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

	if strings.Contains(stringValue, "'") {
		return dataFieldValue, errors.New("DataFieldValueDoesNotAllowSimpleQuote")
	}

	if strings.Contains(stringValue, "\"") {
		return dataFieldValue, errors.New("DataFieldValueDoesNotSupportDoubleQuotes")
	}

	return DataFieldValue(stringValue), nil
}

func (vo DataFieldValue) String() string {
	return string(vo)
}
