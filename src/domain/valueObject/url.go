package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const urlRegex string = `^(?P<schema>https?:\/\/)(?P<hostname>[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9][a-z0-9-]{0,61}[a-z0-9])*)(:(?P<port>\d{1,6}))?(?P<path>\/[A-z0-9\/\_\.\-]*)?(?P<query>\?[\w\/#=&%\-]*)?$`

type Url string

func NewUrl(value string) (Url, error) {
	hasScheme := regexp.MustCompile(`^(http|https)://`)
	if !hasScheme.MatchString(value) {
		value = "https://" + value
	}

	url := Url(value)
	if !url.isValid(value) {
		return "", errors.New("InvalidUrl")
	}
	return url, nil
}

func NewUrlPanic(value string) Url {
	url, err := NewUrl(value)
	if err != nil {
		panic(err)
	}
	return url
}

func (Url) isValid(value string) bool {
	re := regexp.MustCompile(urlRegex)
	return re.MatchString(value)
}

func (url Url) String() string {
	return string(url)
}

func (url Url) getParts() map[string]string {
	return voHelper.FindNamedGroupsMatches(urlRegex, url.String())
}

func (url Url) GetPort() (NetworkPort, error) {
	portStr, exists := url.getParts()["port"]
	if !exists {
		return 0, errors.New("PortNotFound")
	}

	return NewNetworkPort(portStr)
}

func (urlPtr *Url) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	unquotedValue := strings.Trim(valueStr, "\"")

	url, err := NewUrl(unquotedValue)
	if err != nil {
		return err
	}

	*urlPtr = url
	return nil
}
