package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CreateSessionToken struct {
	Username          valueObject.Username         `json:"username"`
	Password          tkValueObject.WeakPassword   `json:"password"`
	OperatorIpAddress tkValueObject.IpAddress      `json:"-"`
}

func NewCreateSessionToken(
	username valueObject.Username,
	password tkValueObject.WeakPassword,
	operatorIpAddress tkValueObject.IpAddress,
) CreateSessionToken {
	return CreateSessionToken{
		Username:          username,
		Password:          password,
		OperatorIpAddress: operatorIpAddress,
	}
}
