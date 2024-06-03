package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const urlPathRegex string = `^(\/|\/\w{1,256}[\w\/\.-]{0,256})$`

type UrlPath string

func NewUrlPath(value interface{}) (urlPath UrlPath, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return urlPath, errors.New("UrlPathValueMustBeString")
	}
	stringValue = strings.TrimSpace(stringValue)

	hasLeadingSlash := strings.HasPrefix(stringValue, "/")
	if !hasLeadingSlash {
		stringValue = "/" + stringValue
	}

	re := regexp.MustCompile(urlPathRegex)
	isValid := re.MatchString(stringValue)
	if !isValid {
		return urlPath, errors.New("InvalidUrlPath")
	}

	return UrlPath(stringValue), nil
}

func NewUrlPathPanic(value interface{}) UrlPath {
	vo, err := NewUrlPath(value)
	if err != nil {
		panic(err)
	}

	return vo
}

func (vo UrlPath) String() string {
	return string(vo)
}

func (vo UrlPath) GetWithoutTrailingSlash() string {
	return strings.TrimSuffix(vo.String(), "/")
}
