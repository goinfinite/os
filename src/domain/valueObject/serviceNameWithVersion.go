package valueObject

import (
	"encoding/json"
	"errors"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type ServiceNameWithVersion struct {
	Name    ServiceName     `json:"name"`
	Version *ServiceVersion `json:"version"`
}

func NewServiceNameWithVersion(
	name ServiceName, version *ServiceVersion,
) ServiceNameWithVersion {
	return ServiceNameWithVersion{
		Name:    name,
		Version: version,
	}
}

func NewServiceNameWithVersionFromString(value interface{}) (
	serviceNameWithVersion ServiceNameWithVersion, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return serviceNameWithVersion, errors.New("ServiceNameWithVersionMustBeString")
	}

	stringValue = strings.ToLower(stringValue)
	stringValueParts := strings.Split(stringValue, ":")
	if len(stringValueParts) == 0 {
		return serviceNameWithVersion, errors.New("InvalidServiceNameWithVersion")
	}

	serviceName, err := NewServiceName(stringValueParts[0])
	if err != nil {
		return serviceNameWithVersion, err
	}

	var serviceVersionPtr *ServiceVersion
	if len(stringValueParts) == 1 {
		return NewServiceNameWithVersion(serviceName, serviceVersionPtr), nil
	}

	serviceVersion, err := NewServiceVersion(stringValueParts[1])
	if err != nil {
		return serviceNameWithVersion, err
	}

	return NewServiceNameWithVersion(serviceName, &serviceVersion), nil
}

func (vo ServiceNameWithVersion) String() string {
	if vo.Version == nil {
		return vo.Name.String()
	}
	return vo.Name.String() + ":" + vo.Version.String()
}

func (vo ServiceNameWithVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(vo.String())
}
