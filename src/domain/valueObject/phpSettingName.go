package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

const phpSettingNameRegex string = `^[a-zA-Z_][a-zA-Z0-9_.-]{1,62}[a-zA-Z0-9_]$`

type PhpSettingName string

func NewPhpSettingName(value interface{}) (settingName PhpSettingName, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
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
