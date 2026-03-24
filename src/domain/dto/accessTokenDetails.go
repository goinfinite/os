package dto

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

type AccessTokenDetails struct {
	TokenType tkValueObject.AccessTokenType `json:"tokenType"`
	AccountId tkValueObject.AccountId       `json:"accountId"`
	IpAddress *tkValueObject.IpAddress      `json:"ipAddress"`
}

func NewAccessTokenDetails(
	tokenType tkValueObject.AccessTokenType,
	accountId tkValueObject.AccountId,
	ipAddress *tkValueObject.IpAddress,
) AccessTokenDetails {
	return AccessTokenDetails{
		TokenType: tokenType,
		AccountId: accountId,
		IpAddress: ipAddress,
	}
}
