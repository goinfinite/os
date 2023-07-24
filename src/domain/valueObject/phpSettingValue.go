package valueObject

import (
	"errors"
)

type PhpSettingValue string

func NewPhpSettingValue(value string) (PhpSettingValue, error) {
	settingValue := PhpSettingValue(value)
	if !settingValue.isValid() {
		return "", errors.New("InvalidPhpSettingValue")
	}
	return settingValue, nil
}

func NewPhpSettingValuePanic(value string) PhpSettingValue {
	settingValue := PhpSettingValue(value)
	if !settingValue.isValid() {
		panic("InvalidPhpSettingValue")
	}
	return settingValue
}

func (settingValue PhpSettingValue) isValid() bool {
	valueLen := len(settingValue)
	return valueLen > 0 && valueLen <= 255
}

func (settingValue PhpSettingValue) String() string {
	return string(settingValue)
}
