package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const serviceEnvRegex string = `^\w{1,1000}=.{1,1000}$`

type ServiceEnv string

func NewServiceEnv(value interface{}) (serviceEnv ServiceEnv, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return serviceEnv, errors.New("ServiceEnvMustBeString")
	}
	stringValue = strings.TrimSpace(stringValue)

	re := regexp.MustCompile(serviceEnvRegex)
	if !re.MatchString(stringValue) {
		return serviceEnv, errors.New("InvalidServiceEnv")
	}

	return ServiceEnv(stringValue), nil
}

func (vo ServiceEnv) String() string {
	return string(vo)
}
