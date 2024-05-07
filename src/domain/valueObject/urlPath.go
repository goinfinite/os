package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

const urlPathRegex string = `^\/[\w/.-]*$`

type UrlPath string

func NewUrlPath(value string) (UrlPath, error) {
	hasLeadingSlash := strings.HasPrefix(value, "/")
	if !hasLeadingSlash {
		value = "/" + value
	}

	compiledRegex := regexp.MustCompile(urlPathRegex)
	isValid := compiledRegex.MatchString(value)
	if !isValid {
		return "", errors.New("InvalidUrlPath")
	}

	return UrlPath(value), nil
}

func NewUrlPathPanic(value string) UrlPath {
	vo, err := NewUrlPath(value)
	if err != nil {
		panic(err)
	}

	return vo
}

func (vo UrlPath) String() string {
	return string(vo)
}
