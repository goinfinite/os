package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UpdatePhpVersionRequest struct {
	Hostname          tkValueObject.Fqdn      `json:"hostname"`
	PhpVersion        valueObject.PhpVersion  `json:"version"`
	OperatorAccountId tkValueObject.AccountId `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress `json:"-"`
}

func NewUpdatePhpVersionRequest(
	hostname tkValueObject.Fqdn,
	phpVersion valueObject.PhpVersion,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) UpdatePhpVersionRequest {
	return UpdatePhpVersionRequest{
		Hostname:          hostname,
		PhpVersion:        phpVersion,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
