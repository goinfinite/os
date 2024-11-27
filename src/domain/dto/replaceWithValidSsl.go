package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReplaceWithValidSsl struct {
	entity.SslPair
	OperatorAccountId valueObject.AccountId `json:"-"`
	OperatorIpAddress valueObject.IpAddress `json:"-"`
}

func NewReplaceWithValidSsl(
	sslPair entity.SslPair,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) ReplaceWithValidSsl {
	return ReplaceWithValidSsl{
		SslPair:           sslPair,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
