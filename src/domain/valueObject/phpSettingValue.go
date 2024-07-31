package valueObject

import (
	"errors"
	"strconv"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type PhpSettingValue string

func NewPhpSettingValue(value interface{}) (settingValue PhpSettingValue, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return settingValue, errors.New("PhpSettingValueMustBeString")
	}
	stringValue = strings.Trim(stringValue, "\"")
	stringValue = strings.ToLower(stringValue)

	if len(stringValue) == 0 {
		return settingValue, errors.New("PhpSettingValueEmpty")
	}

	if len(stringValue) > 255 {
		return settingValue, errors.New("PhpSettingValueTooLong")
	}

	if stringValue == "on" || stringValue == "true" {
		stringValue = "On"
	}

	if stringValue == "off" || stringValue == "false" {
		stringValue = "Off"
	}

	return PhpSettingValue(stringValue), nil
}

func (vo PhpSettingValue) String() string {
	return string(vo)
}

func (vo PhpSettingValue) IsBool() bool {
	return vo == "On" || vo == "Off"
}

func (vo PhpSettingValue) IsNumber() bool {
	_, err := strconv.Atoi(vo.String())
	return err == nil
}

func (vo PhpSettingValue) IsByteSize() bool {
	lastChar := vo[len(vo)-1]
	return lastChar == 'K' || lastChar == 'M' || lastChar == 'G'
}

func (vo PhpSettingValue) GetType() string {
	if vo.IsBool() {
		return "bool"
	}
	if vo.IsNumber() {
		return "number"
	}
	if vo.IsByteSize() {
		return "byteSize"
	}
	return "string"
}
