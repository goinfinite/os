package valueObject

import (
	"errors"
	"regexp"
)

const phpModuleNameRegex string = `^[0-9a-zA-Z_-]{3,32}$`

type PhpModuleName string

func NewPhpModuleName(value string) (PhpModuleName, error) {
	moduleName := PhpModuleName(value)
	if !moduleName.isValid() {
		return "", errors.New("InvalidPhpModuleName")
	}
	return moduleName, nil
}

func NewPhpModuleNamePanic(value string) PhpModuleName {
	moduleName := PhpModuleName(value)
	if !moduleName.isValid() {
		panic("InvalidPhpModuleName")
	}
	return moduleName
}

func (moduleName PhpModuleName) isValid() bool {
	re := regexp.MustCompile(phpModuleNameRegex)
	return re.MatchString(string(moduleName))
}

func (moduleName PhpModuleName) String() string {
	return string(moduleName)
}
