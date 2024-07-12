package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const serviceEnvRegex string = `^\w{1,1000}=.{1,1000}$`

type ServiceEnv string

func NewServiceEnv(value interface{}) (ServiceEnv, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ServiceEnvMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)

	re := regexp.MustCompile(serviceEnvRegex)
	isValid := re.MatchString(stringValue)
	if !isValid {
		return "", errors.New("InvalidServiceEnv")
	}
	return ServiceEnv(stringValue), nil
}

func (vo ServiceEnv) String() string {
	return string(vo)
}
