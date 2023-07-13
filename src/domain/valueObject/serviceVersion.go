package valueObject

import (
	"errors"
	"regexp"
)

const serviceVersionRegex string = `^[\d\.\_\-]{1,20}$`

type ServiceVersion string

func NewServiceVersion(value string) (ServiceVersion, error) {
	user := ServiceVersion(value)
	if !user.isValid() {
		return "", errors.New("InvalidServiceVersion")
	}
	return user, nil
}

func NewServiceVersionPanic(value string) ServiceVersion {
	user := ServiceVersion(value)
	if !user.isValid() {
		panic("InvalidServiceVersion")
	}
	return user
}

func (user ServiceVersion) isValid() bool {
	re := regexp.MustCompile(serviceVersionRegex)
	return re.MatchString(string(user))
}

func (user ServiceVersion) String() string {
	return string(user)
}
