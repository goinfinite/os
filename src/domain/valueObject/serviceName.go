package valueObject

import (
	"errors"
	"regexp"
	"slices"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
	"golang.org/x/exp/maps"
)

type ServiceName string

const ServiceNameRegex string = `^[a-z0-9\.\_\-]{1,64}$`

var NativeSvcNamesWithAliases = map[string][]string{
	"php-webserver": {
		"php", "php-ws", "lsphp", "php-fpm", "php-cgi", "litespeed", "openlitespeed",
	},
	"node": {"nodejs"},
	"mariadb": {
		"mariadbd", "mariadb-server", "mysql", "mysqld", "percona", "perconadb",
	},
	"postgresql": {"postgres"},
	"redis":      {"redis-server"},
	"java-webserver": {
		"java", "jdk", "openjdk", "tomcat",
	},
}

func NewServiceName(value interface{}) (serviceName ServiceName, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return serviceName, errors.New("ServiceNameValueMustBeString")
	}
	svcName := ServiceNameAdapter(stringValue)

	nameRegex := regexp.MustCompile(ServiceNameRegex)
	if !nameRegex.MatchString(svcName) {
		return "", errors.New("InvalidServiceName")
	}

	return ServiceName(svcName), nil
}

func ServiceNameAdapter(value string) string {
	stringValue := strings.TrimSpace(value)
	stringValue = strings.ToLower(stringValue)

	nativeSvcNames := maps.Keys(NativeSvcNamesWithAliases)
	for _, nativeSvcName := range nativeSvcNames {
		if !slices.Contains(NativeSvcNamesWithAliases[nativeSvcName], stringValue) {
			continue
		}
		stringValue = nativeSvcName
		break
	}

	return stringValue
}

func (vo ServiceName) String() string {
	return string(vo)
}
