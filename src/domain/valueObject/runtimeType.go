package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
	"golang.org/x/exp/maps"
)

type RuntimeType string

var runtimeTypesWithAliases = map[string][]string{
	"php-webserver": {
		"php", "php-ws", "lsphp", "php-fpm", "php-cgi", "litespeed", "openlitespeed",
	},
}

func NewRuntimeType(value interface{}) (runtimeType RuntimeType, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return runtimeType, errors.New("RuntimeTypeMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	stringValue, err = runtimeTypeAdapter(stringValue)
	if err != nil {
		return runtimeType, errors.New("InvalidRuntimeType")
	}

	return RuntimeType(stringValue), nil
}

func runtimeTypeAdapter(value string) (string, error) {
	runtimeTypes := maps.Keys(runtimeTypesWithAliases)
	if slices.Contains(runtimeTypes, value) {
		return value, nil
	}

	for _, runtimeType := range runtimeTypes {
		if !slices.Contains(runtimeTypesWithAliases[runtimeType], value) {
			continue
		}

		return runtimeType, nil
	}

	return "", errors.New("InvalidRuntimeType")
}

func (vo RuntimeType) String() string {
	return string(vo)
}
