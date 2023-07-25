package valueObject

import (
	"errors"
	"strconv"
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

func (settingValue PhpSettingValue) IsBool() bool {
	return settingValue == "On" || settingValue == "Off"
}

func (settingValue PhpSettingValue) IsNumber() bool {
	_, err := strconv.Atoi(settingValue.String())
	return err == nil
}

func (settingValue PhpSettingValue) IsByteSize() bool {
	lastChar := settingValue[len(settingValue)-1]
	return lastChar == 'K' || lastChar == 'M' || lastChar == 'G'
}

func (settingValue PhpSettingValue) GetType() string {
	if settingValue.IsBool() {
		return "bool"
	}
	if settingValue.IsNumber() {
		return "number"
	}
	if settingValue.IsByteSize() {
		return "byteSize"
	}
	return "string"
}
