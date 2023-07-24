package valueObject

import (
	"errors"
	"regexp"
)

const phpVersionRegex string = `^\d\.\d$`

type PhpVersion string

func NewPhpVersion(value string) (PhpVersion, error) {
	version := PhpVersion(value)
	if !version.isValid() {
		return "", errors.New("InvalidPhpVersion")
	}
	return version, nil
}

func NewPhpVersionPanic(value string) PhpVersion {
	version := PhpVersion(value)
	if !version.isValid() {
		panic("InvalidPhpVersion")
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
