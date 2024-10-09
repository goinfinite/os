package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type AccessTokenDetails struct {
	TokenType valueObject.AccessTokenType `json:"tokenType"`
	AccountId valueObject.AccountId       `json:"accountId"`
	IpAddress *valueObject.IpAddress      `json:"ipAddress"`
}

func NewAccessTokenDetails(
	tokenType valueObject.AccessTokenType,
	accountId valueObject.AccountId,
	ipAddress *valueObject.IpAddress,
) AccessTokenDetails {
	return AccessTokenDetails{
		TokenType: tokenType,
		AccountId: accountId,
		IpAddress: ipAddress,
	}
}
