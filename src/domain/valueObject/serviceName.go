package valueObject

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type ServiceName string

const ServiceNameRegex string = `^[a-z0-9\.\_\-]{1,64}$`

var NativeSvcNamesWithAliases = map[string][]string{
	"php-webserver": {
		"php",
		"php-ws",
		"lsphp",
		"php-fpm",
		"php-cgi",
		"litespeed",
		"openlitespeed",
	},
	"node": {"nodejs"},
	"mariadb": {
		"mariadbd",
		"mariadb-server",
		"mysql",
		"mysqld",
		"percona",
		"perconadb",
		"percona-server-mysqld",
	},
	"postgresql": {"postgres"},
	"redis":      {"redis-server"},
}

func NewServiceName(value string) (ServiceName, error) {
	svcName := ServiceNameAdapter(value)

	svcNameRegex := regexp.MustCompile(ServiceNameRegex)
	if !svcNameRegex.MatchString(value) {
		return "", errors.New("InvalidServiceName")
	}

	return ServiceName(svcName), nil
}

func NewServiceNamePanic(value string) ServiceName {
	sn, err := NewServiceName(value)
	if err != nil {
		panic(err)
	}
	return sn
}

func ServiceNameAdapter(value string) string {
	svcName := strings.ToLower(value)

	nativeSvcNames := maps.Keys(NativeSvcNamesWithAliases)
	for _, nativeSvcName := range nativeSvcNames {
		if !slices.Contains(NativeSvcNamesWithAliases[nativeSvcName], svcName) {
			continue
		}
		svcName = nativeSvcName
		break
	}

	return svcName
}

func (sn ServiceName) String() string {
	return string(sn)
}
