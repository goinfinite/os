package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var serviceVersionRegex = regexp.MustCompile(
	`^([\w\_\.\-]{1,64}|[\d\.\_\-]{1,20}|latest|lts|alpha|beta|legacy)$`,
)

var serviceVersionPunctuationRegex = regexp.MustCompile(`[\.\_\-]`)

type ServiceVersion string

func NewServiceVersion(value any) (
	serviceVersion ServiceVersion, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return serviceVersion, errors.New("ServiceVersionMustBeString")
	}

	if !serviceVersionRegex.MatchString(stringValue) {
		return serviceVersion, errors.New("InvalidServiceVersion")
	}

	return ServiceVersion(stringValue), nil
}

func (vo ServiceVersion) RemovePunctuation() string {
	return serviceVersionPunctuationRegex.ReplaceAllString(string(vo), "")
}

func (vo ServiceVersion) String() string {
	return string(vo)
}
