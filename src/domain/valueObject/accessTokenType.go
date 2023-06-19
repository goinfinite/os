package valueObject

import "errors"

type AccessTokenType string

const (
	sessionToken AccessTokenType = "sessionToken"
	userApiKey   AccessTokenType = "userApiKey"
)

func NewAccessTokenType(value string) (AccessTokenType, error) {
	att := AccessTokenType(value)
	if !att.isValid() {
		return "", errors.New("InvalidAccessTokenType")
	}
	return att, nil
}

func NewAccessTokenTypePanic(value string) AccessTokenType {
	att := AccessTokenType(value)
	if !att.isValid() {
		panic("InvalidAccessTokenType")
	}
	return att
}

func (att AccessTokenType) isValid() bool {
	switch att {
	case sessionToken, userApiKey:
		return true
	default:
		return false
	}
}

func (att AccessTokenType) String() string {
	return string(att)
}
