package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type ServiceType string

const (
	RuntimeServiceType   ServiceType = "runtime"
	DatabaseServiceType  ServiceType = "database"
	WebServerServiceType ServiceType = "webserver"
	SystemServiceType    ServiceType = "system"
	OtherServiceType     ServiceType = "other"
)

var ValidServiceTypes = []string{
	RuntimeServiceType.String(), DatabaseServiceType.String(),
	WebServerServiceType.String(), SystemServiceType.String(), OtherServiceType.String(),
}

func NewServiceType(value interface{}) (serviceType ServiceType, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return serviceType, errors.New("ServiceTypeValueMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidServiceTypes, stringValue) {
		stringValue = "other"
	}

	return ServiceType(stringValue), nil
}

func (vo ServiceType) String() string {
	return string(vo)
}
