package valueObject

import (
	"errors"
)

type PhpSettingOption string

func NewPhpSettingOption(value string) (PhpSettingOption, error) {
	settingOption := PhpSettingOption(value)
	if !settingOption.isValid() {
		return "", errors.New("InvalidPhpSettingOption")
	}
	return settingOption, nil
}

func NewPhpSettingOptionPanic(value string) PhpSettingOption {
	settingOption := PhpSettingOption(value)
	if !settingOption.isValid() {
		panic("InvalidPhpSettingOption")
	}
	return settingOption
}

func (settingOption PhpSettingOption) isValid() bool {
	valueLen := len(settingOption)
	return valueLen > 0 && valueLen <= 255
}

func (settingOption PhpSettingOption) String() string {
	return string(settingOption)
}
