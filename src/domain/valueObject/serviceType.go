package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type ServiceType string

const (
	ServiceTypeRuntime   ServiceType = "runtime"
	ServiceTypeDatabase  ServiceType = "database"
	ServiceTypeWebServer ServiceType = "webserver"
	ServiceTypeSystem    ServiceType = "system"
	ServiceTypeOther     ServiceType = "other"
)

var ValidServiceTypes = []string{
	ServiceTypeRuntime.String(), ServiceTypeDatabase.String(),
	ServiceTypeWebServer.String(), ServiceTypeSystem.String(),
	ServiceTypeOther.String(),
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
