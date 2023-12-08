package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type ServiceStatus string

var ValidServiceStatuses = []string{
	"running",
	"stopped",
	"uninstalled",
}

func NewServiceStatus(value string) (ServiceStatus, error) {
	value = strings.ToLower(value)
	if !slices.Contains(ValidServiceStatuses, value) {
		return "", errors.New("InvalidServiceStatus")
	}
	return ServiceStatus(value), nil
}

func NewServiceStatusPanic(value string) ServiceStatus {
	ss, err := NewServiceStatus(value)
	if err != nil {
		panic(err)
	}
	return ss
}

func (ss ServiceStatus) String() string {
	return string(ss)
}
