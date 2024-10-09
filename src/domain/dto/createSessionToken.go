package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CreateSessionToken struct {
	Username          valueObject.Username  `json:"username"`
	Password          valueObject.Password  `json:"password"`
	OperatorIpAddress valueObject.IpAddress `json:"-"`
}

func NewCreateSessionToken(
	username valueObject.Username,
	password valueObject.Password,
	operatorIpAddress valueObject.IpAddress,
) CreateSessionToken {
	return CreateSessionToken{
		Username:          username,
		Password:          password,
		OperatorIpAddress: operatorIpAddress,
	}
}
