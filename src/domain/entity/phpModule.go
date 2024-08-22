package entity

import (
	"errors"
	"strings"

	"github.com/speedianet/os/src/domain/valueObject"
	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type PhpModule struct {
	Name   valueObject.PhpModuleName `json:"name"`
	Status bool                      `json:"status"`
}

func NewPhpModule(
	name valueObject.PhpModuleName,
	status bool,
) PhpModule {
	return PhpModule{
		Name:   name,
		Status: status,
	}
}

// format: name:status
func NewPhpModuleFromString(stringValue string) (module PhpModule, err error) {
	stringValueParts := strings.Split(stringValue, ":")
	if len(stringValueParts) == 0 {
		return module, errors.New("EmptyPhpModule")
	}

	if len(stringValueParts) < 2 {
		return module, errors.New("MissingPhpModuleParts")
	}

	name, err := valueObject.NewPhpModuleName(stringValueParts[0])
	if err != nil {
		return module, err
	}

	status, err := voHelper.InterfaceToBool(stringValueParts[1])
	if err != nil {
		return module, err
	}

	return NewPhpModule(name, status), nil
}
