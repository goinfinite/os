package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const serviceVersionRegex string = `^([\w\_\.\-]{1,64}|[\d\.\_\-]{1,20}|latest|lts|alpha|beta)$`

type ServiceVersion string

func NewServiceVersion(value interface{}) (ServiceVersion, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ServiceVersionMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)

	re := regexp.MustCompile(serviceVersionRegex)
	if !re.MatchString(stringValue) {
		return "", errors.New("InvalidServiceVersion")
	}

	return ServiceVersion(stringValue), nil
}

func (vo ServiceVersion) GetWithoutPunctuation() string {
	re := regexp.MustCompile(`[\.\_\-]`)
	return re.ReplaceAllString(string(vo), "")
}

func (vo ServiceVersion) String() string {
	return string(vo)
}
