package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type RuntimeType string

var runtimeTypes = []string{"php"}

func NewRuntimeType(value interface{}) (runtimeType RuntimeType, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return runtimeType, errors.New("RuntimeTypeMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(runtimeTypes, stringValue) {
		return runtimeType, errors.New("InvalidRuntimeType")
	}

	return RuntimeType(stringValue), nil
}

func (vo RuntimeType) String() string {
	return string(vo)
}
