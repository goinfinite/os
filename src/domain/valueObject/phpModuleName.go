package valueObject

import (
	"errors"
	"regexp"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var phpModuleNameRegex = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)

type PhpModuleName string

func NewPhpModuleName(value any) (moduleName PhpModuleName, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return moduleName, errors.New("PhpModuleNameMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !phpModuleNameRegex.MatchString(stringValue) {
		return moduleName, errors.New("InvalidPhpModuleName")
	}

	return PhpModuleName(stringValue), nil
}

func (vo PhpModuleName) String() string {
	return string(vo)
}
