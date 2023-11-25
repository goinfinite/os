package entity

import "github.com/speedianet/os/src/domain/valueObject"

type AccessToken struct {
	Type      valueObject.AccessTokenType `json:"type"`
	ExpiresIn valueObject.UnixTime        `json:"expiresIn"`
	TokenStr  valueObject.AccessTokenStr  `json:"tokenStr"`
}

func NewAccessToken(
	tokenType valueObject.AccessTokenType,
	expiresIn valueObject.UnixTime,
	tokenStr valueObject.AccessTokenStr,
) AccessToken {
	return AccessToken{
		Type:      tokenType,
		ExpiresIn: expiresIn,
		TokenStr:  tokenStr,
	}
}
