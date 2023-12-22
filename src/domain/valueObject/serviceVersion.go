package valueObject

import (
	"errors"
	"regexp"
)

const serviceVersionRegex string = `^([\d\.\_\-]{1,20}|latest|lts|alpha|beta)$`

type ServiceVersion string

func NewServiceVersion(value string) (ServiceVersion, error) {
	version := ServiceVersion(value)
	if !version.isValid() {
		return "", errors.New("InvalidServiceVersion")
	}
	return version, nil
}

func NewServiceVersionPanic(value string) ServiceVersion {
	version, err := NewServiceVersion(value)
	if err != nil {
		panic(err)
	}
	return version
}

func (version ServiceVersion) isValid() bool {
	re := regexp.MustCompile(serviceVersionRegex)
	return re.MatchString(string(version))
}

func (version ServiceVersion) GetWithoutPunctuation() string {
	re := regexp.MustCompile(`[\.\_\-]`)
	return re.ReplaceAllString(string(version), "")
}

func (version ServiceVersion) String() string {
	return string(version)
}
