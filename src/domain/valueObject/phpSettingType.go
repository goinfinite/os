package valueObject

import (
	"errors"
	"slices"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type PhpSettingType string

var validPhpSettingTypes = []string{"select", "text"}

func NewPhpSettingType(value interface{}) (phpSettingType PhpSettingType, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return phpSettingType, errors.New("PhpSettingTypeMustBeString")
	}

	stringValue = strings.ToLower(stringValue)
	if !slices.Contains(validPhpSettingTypes, stringValue) {
		return phpSettingType, errors.New("InvalidPhpSettingType")
	}

	return PhpSettingType(stringValue), nil
}

func (vo PhpSettingType) String() string {
	return string(vo)
}
