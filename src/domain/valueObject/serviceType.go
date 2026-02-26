package valueObject

import (
	"errors"
	"slices"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type ServiceType string

const (
	ServiceTypeSystem    ServiceType = "system"
	ServiceTypeDatabase  ServiceType = "database"
	ServiceTypeRuntime   ServiceType = "runtime"
	ServiceTypeWebServer ServiceType = "webserver"
	ServiceTypeOther     ServiceType = "other"
)

var ValidServiceTypes = []string{
	ServiceTypeSystem.String(), ServiceTypeDatabase.String(),
	ServiceTypeRuntime.String(), ServiceTypeWebServer.String(),
	ServiceTypeOther.String(),
}

func NewServiceType(value interface{}) (serviceType ServiceType, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
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
