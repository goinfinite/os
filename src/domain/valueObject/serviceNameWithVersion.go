package valueObject

import (
	"encoding/json"
	"errors"
	"strings"
)

type ServiceNameWithVersion struct {
	Name    ServiceName     `json:"name"`
	Version *ServiceVersion `json:"version"`
}

func NewServiceNameWithVersion(name ServiceName, version *ServiceVersion) ServiceNameWithVersion {
	return ServiceNameWithVersion{
		Name:    name,
		Version: version,
	}
}

func NewServiceNameWithVersionFromString(value string) (
	serviceNameWithVersion ServiceNameWithVersion, err error,
) {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)
	valueParts := strings.Split(value, ":")
	if len(valueParts) == 0 {
		return serviceNameWithVersion, errors.New("InvalidServiceNameWithVersion")
	}

	serviceName, err := NewServiceName(valueParts[0])
	if err != nil {
		return serviceNameWithVersion, err
	}

	var serviceVersionPtr *ServiceVersion
	if len(valueParts) == 1 {
		return NewServiceNameWithVersion(serviceName, serviceVersionPtr), nil
	}

	serviceVersion, err := NewServiceVersion(valueParts[1])
	if err != nil {
		return serviceNameWithVersion, err
	}

	return NewServiceNameWithVersion(serviceName, &serviceVersion), nil
}

func (vo ServiceNameWithVersion) String() string {
	return vo.Name.String() + ":" + vo.Version.String()
}

func (vo ServiceNameWithVersion) MarshalJSON() ([]byte, error) {
	if vo.Version == nil {
		return json.Marshal(vo.Name.String())
	}
	return json.Marshal(vo.String())
}
