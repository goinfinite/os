package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type ServiceType string

const (
	ServiceTypeDatabase  ServiceType = "database"
	ServiceTypeOther     ServiceType = "other"
	ServiceTypeRuntime   ServiceType = "runtime"
	ServiceTypeSystem    ServiceType = "system"
	ServiceTypeWebServer ServiceType = "webserver"
)

var ValidServiceTypes = []string{
	ServiceTypeDatabase.String(), ServiceTypeOther.String(),
	ServiceTypeRuntime.String(), ServiceTypeOther.String(),
	ServiceTypeWebServer.String(),
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
