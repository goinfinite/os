package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var dataFieldNameRegex = regexp.MustCompile(`^\w[\w-]{1,128}\w$`)

type DataFieldName string

func NewDataFieldName(value any) (
	dataFieldName DataFieldName, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return dataFieldName, errors.New("DataFieldNameMustBeString")
	}

	if !dataFieldNameRegex.MatchString(stringValue) {
		return dataFieldName, errors.New("InvalidDataFieldName")
	}

	return DataFieldName(stringValue), nil
}

func (vo DataFieldName) String() string {
	return string(vo)
}
