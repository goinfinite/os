package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type ServiceDescription string

func NewServiceDescription(value interface{}) (ServiceDescription, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ServiceDescriptionValueMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)

	if len(stringValue) < 2 {
		return "", errors.New("ServiceDescriptionTooSmall")
	}

	if len(stringValue) > 2048 {
		return "", errors.New("ServiceDescriptionTooBig")
	}

	return ServiceDescription(stringValue), nil
}

func (vo ServiceDescription) String() string {
	return string(vo)
}
