package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteSslPair struct {
	SslPairId         valueObject.SslPairId `json:"sslPairId"`
	OperatorAccountId valueObject.AccountId `json:"-"`
	OperatorIpAddress valueObject.IpAddress `json:"-"`
}

func NewDeleteSslPair(
	sslPairId valueObject.SslPairId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteSslPair {
	return DeleteSslPair{
		SslPairId:         sslPairId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
