package valueObject

import (
	"errors"

	"golang.org/x/exp/slices"
)

type ServiceName string

var SupportedServiceNames = []string{
	"openlitespeed",
	"nginx",
	"node",
	"mysql",
	"redis",
}

var SupportedServiceNamesAliases = []string{
	"litespeed",
	"nodejs",
	"mysqld",
	"mariadb",
	"percona",
	"perconadb",
	"redis-server",
}

func NewServiceName(value string) (ServiceName, error) {
	ss := ServiceName(value)
	if !ss.isValid() {
		return "", errors.New("InvalidServiceName")
	}
	return ss, nil
}

func NewServiceNamePanic(value string) ServiceName {
	ss := ServiceName(value)
	if !ss.isValid() {
		panic("InvalidServiceName")
	}
	return ss
}

func (ss ServiceName) isValid() bool {
	supportedServices := append(SupportedServiceNames, SupportedServiceNamesAliases...)
	return slices.Contains(supportedServices, ss.String())
}

func (ss ServiceName) String() string {
	return string(ss)
}
