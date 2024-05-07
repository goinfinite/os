package valueObject

import (
	"errors"
	"slices"
	"strings"
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
	sn, err := NewServiceNature(value)
	if err != nil {
		panic(err)
	}
	return sn
}

func (sn ServiceNature) String() string {
	return string(sn)
}
