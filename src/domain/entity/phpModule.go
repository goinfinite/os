package entity

import "github.com/speedianet/os/src/domain/valueObject"

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
