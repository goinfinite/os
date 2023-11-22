package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type ServiceType string

var ValidServiceTypes = []string{
	"runtime",
	"database",
}

func NewServiceType(value string) (ServiceType, error) {
	st := ServiceType(strings.ToLower(value))
	if !st.isValid() {
		return "", errors.New("InvalidServiceType")
	}
	return st, nil
}

func NewServiceTypePanic(value string) ServiceType {
	st, err := NewServiceType(value)
	if err != nil {
		panic(err)
	}
	return st
}

func (st ServiceType) isValid() bool {
	return slices.Contains(ValidServiceTypes, st.String())
}

func (st ServiceType) String() string {
	return string(st)
}
