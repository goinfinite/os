package valueObject

import (
	"errors"
	"regexp"
)

const phpSettingNameRegex string = `^\w[\w\.-]{1,62}\w$`

type PhpSettingName string

func NewPhpSettingName(value string) (PhpSettingName, error) {
	settingName := PhpSettingName(value)
	if !settingName.isValid() {
		return "", errors.New("InvalidPhpSettingName")
	}
	return settingName, nil
}

func NewPhpSettingNamePanic(value string) PhpSettingName {
	settingName := PhpSettingName(value)
	if !settingName.isValid() {
		panic("InvalidPhpSettingName")
	}
	return settingName
}

func (settingName PhpSettingName) isValid() bool {
	re := regexp.MustCompile(phpSettingNameRegex)
	return re.MatchString(string(settingName))
}

func (settingName PhpSettingName) String() string {
	return string(settingName)
}
