package valueObject

import (
	"errors"
)

type DataFieldValue string

func NewDataFieldValue(value string) (DataFieldValue, error) {
	if len(value) <= 1 {
		return "", errors.New("DataFieldValueTooSmall")
	}

	if len(value) >= 60 {
		return "", errors.New("DataFieldValueTooBig")
	}

	return DataFieldValue(value), nil
}

func NewDataFieldValuePanic(value string) DataFieldValue {
	dfv, err := NewDataFieldValue(value)
	if err != nil {
		panic(err)
	}

	return dfv
}

func (dfv DataFieldValue) String() string {
	return string(dfv)
}
