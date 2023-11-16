package valueObject

import (
	"errors"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type ServiceName string

var SupportedServiceNamesAndAliases = map[string][]string{
	"openlitespeed": {"litespeed"},
	"node":          {"nodejs"},
	"mysql":         {"mysqld", "mariadb", "percona", "perconadb"},
	"redis":         {"redis-server"},
	"php":           {"lsphp", "php-fpm", "php-cgi", "php-cli"},
}

func NewServiceName(value string) (ServiceName, error) {
	servicesName := maps.Keys(SupportedServiceNamesAndAliases)
	if slices.Contains(servicesName, value) {
		return ServiceName(value), nil
	}

	for _, serviceName := range servicesName {
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
	sn, err := NewServiceName(value)
	if err != nil {
		panic(err)
	}
	return sn
}

func (sn ServiceName) String() string {
	return string(sn)
}
