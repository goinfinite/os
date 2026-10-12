package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var serviceEnvRegex = regexp.MustCompile(`^\w{1,1000}=.{1,1000}$`)

type ServiceEnv string

func NewServiceEnv(value any) (serviceEnv ServiceEnv, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return serviceEnv, errors.New("ServiceEnvMustBeString")
	}

	if !serviceEnvRegex.MatchString(stringValue) {
		return serviceEnv, errors.New("InvalidServiceEnv")
	}

	return ServiceEnv(stringValue), nil
}

func (vo ServiceEnv) String() string {
	return string(vo)
}
