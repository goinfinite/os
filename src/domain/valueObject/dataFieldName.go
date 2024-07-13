package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const dataFieldNameRegex string = `^\w[\w-]{1,128}\w$`

type DataFieldName string

func NewDataFieldName(value interface{}) (DataFieldName, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("DataFieldNameMustBeString")
	}

	re := regexp.MustCompile(dataFieldNameRegex)
	if !re.MatchString(stringValue) {
		return "", errors.New("InvalidDataFieldName")
	}

	return DataFieldName(stringValue), nil
}

func NewDataFieldNamePanic(value interface{}) DataFieldName {
	dfn, err := NewDataFieldName(value)
	if err != nil {
		panic(err)
	}

	return dfn
}

func (vo DataFieldName) String() string {
	return string(vo)
}
