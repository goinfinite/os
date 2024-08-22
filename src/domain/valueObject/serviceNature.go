package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type ServiceNature string

var ValidServiceNatures = []string{
	"solo", "multi", "custom",
}

func NewServiceNature(value interface{}) (serviceNature ServiceNature, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return serviceNature, errors.New("ServiceNatureValueMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidServiceNatures, stringValue) {
		return serviceNature, errors.New("InvalidServiceNature")
	}

	return ServiceNature(stringValue), nil
}

func (vo ServiceNature) String() string {
	return string(vo)
}
