package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type ServiceNature string

var ValidServiceNatures = []string{
	"solo",
	"multi",
	"custom",
}

func NewServiceNature(value string) (ServiceNature, error) {
	value = strings.ToLower(value)
	if !slices.Contains(ValidServiceNatures, value) {
		return "", errors.New("InvalidServiceNature")
	}
	return ServiceNature(value), nil
}

func NewServiceNaturePanic(value string) ServiceNature {
	st, err := NewServiceNature(value)
	if err != nil {
		panic(err)
	}
	return st
}

func (st ServiceNature) String() string {
	return string(st)
}
