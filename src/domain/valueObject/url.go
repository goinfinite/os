package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const urlRegexExpression string = `^(?P<schema>https?:\/\/)(?P<hostname>[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9][a-z0-9-]{0,61}[a-z0-9])*)(:(?P<port>\d{1,6}))?(?P<path>\/[A-Za-z0-9\/\_\.\-]*)?(?P<query>\?[\w\/#=&%\-]*)?$`

type Url string

func NewUrl(value interface{}) (url Url, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return url, errors.New("UrlValueMustBeString")
	}

	if !strings.HasPrefix(stringValue, "http") {
		stringValue = "https://" + stringValue
	}

	urlRegex := regexp.MustCompile(urlRegexExpression)
	if !urlRegex.MatchString(stringValue) {
		return url, errors.New("InvalidUrl")
	}

	return Url(stringValue), nil
}

func (vo Url) String() string {
	return string(vo)
}
