package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

const unixFileExtensionRegexExpression = `^[\w\_\-]{1,15}$`

type UnixFileExtension string

func NewUnixFileExtension(value string) (UnixFileExtension, error) {
	if strings.HasPrefix(value, ".") {
		value, _ = strings.CutPrefix(value, ".")
	}

	unixFileExtension := UnixFileExtension(value)
	if !unixFileExtension.isValid() {
		return "", errors.New("InvalidUnixFileExtension")
	}
	return unixFileExtension, nil
}

func NewUnixFileExtensionPanic(value string) UnixFileExtension {
	unixFileExtension, err := NewUnixFileExtension(value)
	if err != nil {
		panic(err)
	}
	return unixFileExtension
}

func (unixFileExtension UnixFileExtension) isValid() bool {
	unixFileExtensionRegex := regexp.MustCompile(unixFileExtensionRegexExpression)
	isValidFormat := unixFileExtensionRegex.MatchString(string(unixFileExtension))

	isEmpty := string(unixFileExtension) == ""

	return isValidFormat && !isEmpty
}

func (unixFileExtension UnixFileExtension) String() string {
	return string(unixFileExtension)
}
