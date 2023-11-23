package valueObject

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

const ServiceNameRegex string = `^[a-z0-9\.\_\-]{1,64}$`

type ServiceName string

var NativeSvcNamesWithAliases = map[string][]string{
	"php":   {"lsphp", "php-fpm", "php-cgi", "litespeed", "openlitespeed"},
	"node":  {"nodejs"},
	"mysql": {"mysqld", "mariadb", "percona", "perconadb"},
	"redis": {"redis-server"},
}

func NewServiceName(value string) (ServiceName, error) {
	value = strings.ToLower(value)

	servicesName := maps.Keys(NativeSvcNamesWithAliases)
	for _, serviceName := range servicesName {
		if !slices.Contains(NativeSvcNamesWithAliases[serviceName], value) {
			continue
		}
		value = serviceName
	}

	svcName := ServiceName(value)
	if !svcName.isValid() {
		return "", errors.New("InvalidServiceName")
	}

	return svcName, nil
}

func NewServiceNamePanic(value string) ServiceName {
	sn, err := NewServiceName(value)
	if err != nil {
		panic(err)
	}
	return sn
}

func (sn ServiceName) isValid() bool {
	re := regexp.MustCompile(ServiceNameRegex)
	return re.MatchString(string(sn))
}

func (sn ServiceName) String() string {
	return string(sn)
}
