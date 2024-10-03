package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type PhpSettingType string

var validPhpSettingTypes = []string{"select", "text"}

func NewPhpSettingType(value interface{}) (phpSettingType PhpSettingType, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
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
