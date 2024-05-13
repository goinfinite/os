package valueObject

import (
	"errors"
	"regexp"
)

const dataFieldNameRegex string = `^\w[\w-]{1,128}\w$`

type DataFieldName string

func NewDataFieldName(value string) (DataFieldName, error) {
	dfn := DataFieldName(value)
	if !dfn.isValid() {
		return "", errors.New("InvalidDataFieldName")
	}

	return dfn, nil
}

func NewDataFieldNamePanic(value string) DataFieldName {
	dfn, err := NewDataFieldName(value)
	if err != nil {
		panic(err)
	}

	return dfn
}

func (dfn DataFieldName) isValid() bool {
	re := regexp.MustCompile(dataFieldNameRegex)
	return re.MatchString(string(dfn))
}

func (dfn DataFieldName) String() string {
	return string(dfn)
}
