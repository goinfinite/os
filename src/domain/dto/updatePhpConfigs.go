package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type UpdatePhpConfigs struct {
	Hostname    valueObject.Fqdn       `json:"hostname"`
	PhpVersion  valueObject.PhpVersion `json:"version"`
	PhpModules  []entity.PhpModule     `json:"modules"`
	PhpSettings []entity.PhpSetting    `json:"settings"`
}

func NewUpdatePhpConfigs(
	hostname valueObject.Fqdn,
	phpVersion valueObject.PhpVersion,
	phpModules []entity.PhpModule,
	phpSettings []entity.PhpSetting,
) UpdatePhpConfigs {
	return UpdatePhpConfigs{
		Hostname:    hostname,
		PhpVersion:  phpVersion,
		PhpModules:  phpModules,
		PhpSettings: phpSettings,
	}
}
