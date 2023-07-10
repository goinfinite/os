package valueObject

import (
	"errors"
	"regexp"
)

const serviceNameRegex string = `^[a-z_]([a-z0-9_-]{0,31}|[a-z0-9_-]{0,30}\$)$`

type ServiceName string

func NewServiceName(value string) (ServiceName, error) {
	user := ServiceName(value)
	if !user.isValid() {
		return "", errors.New("InvalidServiceName")
	}
	return user, nil
}

func NewServiceNamePanic(value string) ServiceName {
	user := ServiceName(value)
	if !user.isValid() {
		panic("InvalidServiceName")
	}
	return user
}

func (user ServiceName) isValid() bool {
	re := regexp.MustCompile(serviceNameRegex)
	return re.MatchString(string(user))
}

func (user ServiceName) String() string {
	return string(user)
}
