package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var dataFieldLabelRegex = regexp.MustCompile(`^\w[\w- ]{1,127}\w$`)

type DataFieldLabel string

func NewDataFieldLabel(value any) (
	dataFieldLabel DataFieldLabel, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return dataFieldLabel, errors.New("DataFieldLabelMustBeString")
	}

	if !dataFieldLabelRegex.MatchString(stringValue) {
		return dataFieldLabel, errors.New("InvalidDataFieldLabel")
	}

	return DataFieldLabel(stringValue), nil
}

func (vo DataFieldLabel) String() string {
	return string(vo)
}
