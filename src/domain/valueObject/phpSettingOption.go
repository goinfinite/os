package valueObject

import (
	"errors"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type PhpSettingOption string

func NewPhpSettingOption(value interface{}) (
	phpSettingOption PhpSettingOption, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return phpSettingOption, errors.New("PhpSettingOptionMustBeString")
	}

	if len(stringValue) == 0 {
		return phpSettingOption, errors.New("EmptyPhpSettingOption")
	}

	if len(stringValue) > 255 {
		return phpSettingOption, errors.New("PhpSettingOptionTooLong")
	}

	return PhpSettingOption(stringValue), nil
}

func (vo PhpSettingOption) String() string {
	return string(vo)
}
