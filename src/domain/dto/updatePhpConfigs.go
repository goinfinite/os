package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type UpdatePhpConfigs struct {
	Hostname          valueObject.Fqdn       `json:"hostname"`
	PhpVersion        valueObject.PhpVersion `json:"version"`
	PhpModules        []entity.PhpModule     `json:"modules"`
	PhpSettings       []entity.PhpSetting    `json:"settings"`
	OperatorAccountId valueObject.AccountId  `json:"-"`
	OperatorIpAddress valueObject.IpAddress  `json:"-"`
}

func NewUpdatePhpConfigs(
	hostname valueObject.Fqdn,
	phpVersion valueObject.PhpVersion,
	phpModules []entity.PhpModule,
	phpSettings []entity.PhpSetting,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) UpdatePhpConfigs {
	return UpdatePhpConfigs{
		Hostname:          hostname,
		PhpVersion:        phpVersion,
		PhpModules:        phpModules,
		PhpSettings:       phpSettings,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
