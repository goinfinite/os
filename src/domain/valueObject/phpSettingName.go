package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const phpSettingNameRegex string = `^\w[\w\.-]{1,62}\w$`

type PhpSettingName string

func NewPhpSettingName(value interface{}) (settingName PhpSettingName, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return settingName, errors.New("PhpSettingNameMustBeString")
	}

	re := regexp.MustCompile(phpSettingNameRegex)
	if !re.MatchString(stringValue) {
		return settingName, errors.New("InvalidPhpSettingName")
	}

	return PhpSettingName(stringValue), nil
}

func (vo PhpSettingName) String() string {
	return string(vo)
}
