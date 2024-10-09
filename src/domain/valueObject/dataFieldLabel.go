package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const dataFieldLabelRegex string = `^\w[\w- ]{1,127}\w$`

type DataFieldLabel string

func NewDataFieldLabel(value interface{}) (
	dataFieldLabel DataFieldLabel, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return dataFieldLabel, errors.New("DataFieldLabelMustBeString")
	}

	re := regexp.MustCompile(dataFieldLabelRegex)
	if !re.MatchString(stringValue) {
		return dataFieldLabel, errors.New("InvalidDataFieldLabel")
	}

	return DataFieldLabel(stringValue), nil
}

func (vo DataFieldLabel) String() string {
	return string(vo)
}
