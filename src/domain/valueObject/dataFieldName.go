package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

const dataFieldNameRegex string = `^\w[\w-]{1,128}\w$`

type DataFieldName string

func NewDataFieldName(value interface{}) (
	dataFieldName DataFieldName, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return dataFieldName, errors.New("DataFieldNameMustBeString")
	}

	re := regexp.MustCompile(dataFieldNameRegex)
	if !re.MatchString(stringValue) {
		return dataFieldName, errors.New("InvalidDataFieldName")
	}

	return DataFieldName(stringValue), nil
}

func (vo DataFieldName) String() string {
	return string(vo)
}
