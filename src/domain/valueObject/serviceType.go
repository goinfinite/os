package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type ServiceType string

var ValidServiceTypes = []string{
	"application", "runtime", "database", "webserver", "mom", "monitoring",
	"logging", "security", "backup", "system", "other",
}

func NewServiceType(value interface{}) (serviceType ServiceType, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return serviceType, errors.New("ServiceTypeValueMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidServiceTypes, stringValue) {
		return serviceType, errors.New("InvalidServiceType")
	}

	return ServiceType(stringValue), nil
}

func (vo ServiceType) String() string {
	return string(vo)
}
