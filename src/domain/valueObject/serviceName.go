package valueObject

import (
	"errors"

	"golang.org/x/exp/slices"
)

type ServiceName string

var SupportedServiceNames = []string{
	"openlitespeed",
	"litespeed",
	"nginx",
	"node",
	"nodejs",
	"mysql",
	"mysqld",
	"mariadb",
	"percona",
	"perconadb",
	"redis",
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
	return slices.Contains(SupportedServiceNames, ss.String())
}

func (ss ServiceName) String() string {
	return string(ss)
}
