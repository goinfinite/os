package dto

import "github.com/speedianet/sam/src/domain/valueObject"

type AccessTokenDetails struct {
	TokenType valueObject.AccessTokenType `json:"tokenType"`
	UserId    valueObject.UserId          `json:"userId"`
	IpAddress *valueObject.IpAddress      `json:"ipAddress"`
}

func NewAccessTokenDetails(
	tokenType valueObject.AccessTokenType,
	userId valueObject.UserId,
	ipAddress *valueObject.IpAddress,
) AccessTokenDetails {
	return AccessTokenDetails{
		TokenType: tokenType,
		UserId:    userId,
		IpAddress: ipAddress,
	}
}
