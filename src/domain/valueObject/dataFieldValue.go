package valueObject

import (
	"errors"
	"strings"
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

func (dfvPtr *DataFieldValue) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	unquotedValue := strings.Trim(valueStr, "\"")

	dfv, err := NewDataFieldValue(unquotedValue)
	if err != nil {
		return err
	}

	*dfvPtr = dfv
	return nil
}
