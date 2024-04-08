package valueObject

import (
	"errors"
	"regexp"
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
