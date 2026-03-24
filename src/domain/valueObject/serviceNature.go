package valueObject

import (
	"errors"
	"slices"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type ServiceNature string

const (
	ServiceNatureSolo   ServiceNature = "solo"
	ServiceNatureMulti  ServiceNature = "multi"
	ServiceNatureCustom ServiceNature = "custom"
)

var ValidServiceNatures = []string{
	ServiceNatureSolo.String(), ServiceNatureMulti.String(),
	ServiceNatureCustom.String(),
}

func NewServiceNature(value interface{}) (serviceNature ServiceNature, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
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
