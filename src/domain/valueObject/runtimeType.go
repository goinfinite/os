package valueObject

import (
	"errors"
	"slices"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type RuntimeType string

var runtimeTypes = []string{"php"}

func NewRuntimeType(value interface{}) (runtimeType RuntimeType, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
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
