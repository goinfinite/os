package valueObject

import (
	"errors"
	"slices"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type AccessTokenType string

var validAccessTokenType = []string{
	"sessionToken", "accountApiKey",
}

func NewAccessTokenType(value interface{}) (
	accessTokenType AccessTokenType, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return accessTokenType, errors.New("AccessTokenTypeMustBeString")
	}

	if !slices.Contains(validAccessTokenType, stringValue) {
		return accessTokenType, errors.New("InvalidAccessTokenType")
	}

	return AccessTokenType(stringValue), nil
}

func (vo AccessTokenType) String() string {
	return string(vo)
}
