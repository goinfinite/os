package valueObject

import (
	"errors"
	"regexp"
)

const phpVersionRegex string = `^\d\.\d$`

type PhpVersion string

func NewPhpVersion(value string) (PhpVersion, error) {
	version := PhpVersion(value)
	if len(value) == 2 {
		version = PhpVersion(value[:1] + "." + value[1:])
	}

	if !version.isValid() {
		return "", errors.New("InvalidPhpVersion")
	}
	return version, nil
}

func NewPhpVersionPanic(value string) PhpVersion {
	version, err := NewPhpVersion(value)
	if err != nil {
		panic(err)
	}

	return version
}

func (version PhpVersion) isValid() bool {
	re := regexp.MustCompile(phpVersionRegex)
	return re.MatchString(string(version))
}

func (version PhpVersion) String() string {
	return string(version)
}

func (version PhpVersion) GetWithoutDots() string {
	return string(version[:1]) + string(version[2:])
}
