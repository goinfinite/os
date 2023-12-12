package valueObject

import (
	"errors"
	"regexp"
)

const urlPathRegex string = `^(?P<path>\/[a-z0-9\/\_\.\-]*)?(?P<query>\?[\w\/#=&%\-]*)?$`

type UrlPath string

func NewUrlPath(value string) (UrlPath, error) {
	urlPath := UrlPath(value)
	if !urlPath.isValid(value) {
		return "", errors.New("InvalidUrlPath")
	}
	return urlPath, nil
}

func NewUrlPathPanic(value string) UrlPath {
	urlPath, err := NewUrlPath(value)
	if err != nil {
		panic(err)
	}
	return urlPath
}

func (UrlPath) isValid(value string) bool {
	re := regexp.MustCompile(urlPathRegex)
	return re.MatchString(value)
}

func (urlPath UrlPath) String() string {
	return string(urlPath)
}
