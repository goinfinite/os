package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UpdatePhpConfigs struct {
	Hostname          tkValueObject.Fqdn     `json:"hostname"`
	PhpVersion        valueObject.PhpVersion `json:"version"`
	PhpModules        []entity.PhpModule     `json:"modules"`
	PhpSettings       []entity.PhpSetting    `json:"settings"`
	OperatorAccountId tkValueObject.AccountId `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress `json:"-"`
}

func NewUpdatePhpConfigs(
	hostname tkValueObject.Fqdn,
	phpVersion valueObject.PhpVersion,
	phpModules []entity.PhpModule,
	phpSettings []entity.PhpSetting,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
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
