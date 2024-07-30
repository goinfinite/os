package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const phpVersionRegex string = `^\d\.\d$`

type PhpVersion string

func NewPhpVersion(value interface{}) (phpVersion PhpVersion, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return phpVersion, errors.New("PhpVersionMustBeString")
	}

	if len(stringValue) == 2 {
		stringValue = stringValue[:1] + "." + stringValue[1:]
	}

	re := regexp.MustCompile(phpVersionRegex)
	isValid := re.MatchString(stringValue)
	if !isValid {
		return "", errors.New("InvalidPhpVersion")
	}
	return PhpVersion(stringValue), nil
}

func (vo PhpVersion) String() string {
	return string(vo)
}

func (vo PhpVersion) GetWithoutDots() string {
	return string(vo[:1]) + string(vo[2:])
}
