package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

const dataFieldKeyRegex string = `^[0-9a-zA-Z_-]{1,32}$`

type DataFieldKey string

func NewDataFieldKey(value string) (DataFieldKey, error) {
	dfk := DataFieldKey(value)
	if !dfk.isValid() {
		return "", errors.New("InvalidDataFieldKey")
	}

	return dfk, nil
}

func NewDataFieldKeyPanic(value string) DataFieldKey {
	dfk, err := NewDataFieldKey(value)
	if err != nil {
		panic(err)
	}

	return dfk
}

func (dfk DataFieldKey) isValid() bool {
	re := regexp.MustCompile(dataFieldKeyRegex)
	return re.MatchString(string(dfk))
}

func (dfk DataFieldKey) String() string {
	return string(dfk)
}

func (dfkPtr *DataFieldKey) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	unquotedValue := strings.Trim(valueStr, "\"")

	dfk, err := NewDataFieldKey(unquotedValue)
	if err != nil {
		return err
	}

	*dfkPtr = dfk
	return nil
}
