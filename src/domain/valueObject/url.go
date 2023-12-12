package valueObject

import (
	"errors"
	"regexp"
)

const urlRegex string = `^(?P<schema>https?:\/\/)(?P<fqdn>[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9][a-z0-9-]{0,61}[a-z0-9])+)(?P<port>:\d{1,6})?(?P<path>\/[a-z0-9\/\_\.\-]*)?(?P<query>\?[\w\/#=&%\-]*)?$`

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
