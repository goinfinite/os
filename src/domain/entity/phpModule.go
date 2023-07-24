package entity

import "github.com/speedianet/sam/src/domain/valueObject"

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
