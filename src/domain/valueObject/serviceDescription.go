package valueObject

import (
	"errors"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type ServiceDescription string

func NewServiceDescription(value interface{}) (
	serviceDescription ServiceDescription, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return serviceDescription, errors.New("ServiceDescriptionValueMustBeString")
	}

	if len(stringValue) < 2 {
		return serviceDescription, errors.New("ServiceDescriptionTooSmall")
	}

	if len(stringValue) > 2048 {
		return serviceDescription, errors.New("ServiceDescriptionTooBig")
	}

	return ServiceDescription(stringValue), nil
}

func (vo ServiceDescription) String() string {
	return string(vo)
}
