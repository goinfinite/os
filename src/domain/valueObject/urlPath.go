package valueObject

import (
	"errors"
	"regexp"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

const urlPathRegex string = `^(\/|\/\w{1,256}[\w\/\.-]{0,256})$`

type UrlPath string

func NewUrlPath(value interface{}) (urlPath UrlPath, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return urlPath, errors.New("UrlPathValueMustBeString")
	}

	hasLeadingSlash := strings.HasPrefix(stringValue, "/")
	if !hasLeadingSlash {
		stringValue = "/" + stringValue
	}

	re := regexp.MustCompile(urlPathRegex)
	if !re.MatchString(stringValue) {
		return urlPath, errors.New("InvalidUrlPath")
	}

	return UrlPath(stringValue), nil
}

func (vo UrlPath) String() string {
	return string(vo)
}

func (vo UrlPath) GetWithoutTrailingSlash() string {
	return strings.TrimSuffix(vo.String(), "/")
}
