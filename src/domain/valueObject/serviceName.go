package valueObject

import (
	"errors"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type ServiceName string

var SupportedServiceNamesAndAliases = map[string][]string{
	"openlitespeed": {"litespeed"},
	"nginx":         {},
	"node":          {"nodejs"},
	"mysql":         {"mysqld", "mariadb", "percona", "perconadb"},
	"redis":         {"redis-server"},
}

func NewServiceName(value string) (ServiceName, error) {
	supportedServicesCorrectName := maps.Keys(SupportedServiceNamesAndAliases)
	if slices.Contains(supportedServicesCorrectName, value) {
		return ServiceName(value), nil
	}

	for _, serviceName := range supportedServicesCorrectName {
		if slices.Contains(
			SupportedServiceNamesAndAliases[serviceName],
			value,
		) {
			return ServiceName(value), nil
		}
	}

	return "", errors.New("InvalidServiceName")
}

func NewServiceNamePanic(value string) ServiceName {
	ss, err := NewServiceName(value)
	if err != nil {
		panic("InvalidServiceName")
	}
	return ss
}

func (ss ServiceName) String() string {
	return string(ss)
}
