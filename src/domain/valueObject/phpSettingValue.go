package valueObject

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type PhpSettingValue string

const (
	PhpSettingValueTypeBool     string = "bool"
	PhpSettingValueTypeNumber   string = "number"
	PhpSettingValueTypeByteSize string = "byteSize"
	PhpSettingValueTypeString   string = "string"
)

func NewPhpSettingValue(value any) (settingValue PhpSettingValue, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return settingValue, errors.New("PhpSettingValueMustBeString")
	}
	stringValue = strings.Trim(stringValue, "\"")

	if len(stringValue) == 0 {
		return settingValue, errors.New("EmptyPhpSettingValue")
	}

	if len(stringValue) > 255 {
		return settingValue, errors.New("PhpSettingValueTooLong")
	}

	hasControlChars := strings.IndexFunc(stringValue, unicode.IsControl) > -1
	if hasControlChars {
		return settingValue, errors.New("PhpSettingValueHasControlChars")
	}

	switch strings.ToLower(stringValue) {
	case "on", "true":
		stringValue = "On"
	case "off", "false":
		stringValue = "Off"
	}

	settingValue = PhpSettingValue(stringValue)
	if settingValue.IsByteSize() {
		settingValue = PhpSettingValue(strings.ToUpper(settingValue.String()))
	}

	return settingValue, nil
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
	if len(vo) < 2 {
		return false
	}

	lastChar := vo[len(vo)-1]
	isByteUnit := strings.IndexByte("KkMmGg", lastChar) > -1
	if !isByteUnit {
		return false
	}

	amount := string(vo[:len(vo)-1])
	_, err := strconv.ParseUint(amount, 10, 64)
	return err == nil
}

func (vo PhpSettingValue) ReadType() string {
	if vo.IsBool() {
		return PhpSettingValueTypeBool
	}
	if vo.IsNumber() {
		return PhpSettingValueTypeNumber
	}
	if vo.IsByteSize() {
		return PhpSettingValueTypeByteSize
	}
	return PhpSettingValueTypeString
}
