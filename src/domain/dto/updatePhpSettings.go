package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UpdatePhpSettingsRequest struct {
	Hostname          tkValueObject.Fqdn      `json:"hostname"`
	PhpSettings       []entity.PhpSetting     `json:"settings"`
	OperatorAccountId tkValueObject.AccountId `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress `json:"-"`
}

func NewUpdatePhpSettingsRequest(
	hostname tkValueObject.Fqdn,
	phpSettings []entity.PhpSetting,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) UpdatePhpSettingsRequest {
	return UpdatePhpSettingsRequest{
		Hostname:          hostname,
		PhpSettings:       phpSettings,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
