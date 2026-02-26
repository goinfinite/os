package entity

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

type AccessToken struct {
	Type      tkValueObject.AccessTokenType  `json:"type"`
	ExpiresIn tkValueObject.UnixTime         `json:"expiresIn"`
	TokenStr  tkValueObject.AccessTokenValue `json:"tokenStr"`
}

func NewAccessToken(
	tokenType tkValueObject.AccessTokenType,
	expiresIn tkValueObject.UnixTime,
	tokenStr tkValueObject.AccessTokenValue,
) AccessToken {
	return AccessToken{
		Type:      tokenType,
		ExpiresIn: expiresIn,
		TokenStr:  tokenStr,
	}
}
