package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var phpVersionRegex = regexp.MustCompile(`^\d\.\d$`)

type PhpVersion string

func NewPhpVersion(value any) (phpVersion PhpVersion, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return phpVersion, errors.New("PhpVersionMustBeString")
	}

	if len(stringValue) == 2 {
		stringValue = stringValue[:1] + "." + stringValue[1:]
	}

	if !phpVersionRegex.MatchString(stringValue) {
		return "", errors.New("InvalidPhpVersion")
	}

	return PhpVersion(stringValue), nil
}

func (vo PhpVersion) String() string {
	return string(vo)
}

func (vo PhpVersion) RemoveDots() string {
	return string(vo[:1]) + string(vo[2:])
}
