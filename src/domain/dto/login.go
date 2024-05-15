package dto

import "github.com/speedianet/os/src/domain/valueObject"

type Login struct {
	Username  valueObject.Username  `json:"username"`
	Password  valueObject.Password  `json:"password"`
	IpAddress valueObject.IpAddress `json:"-"`
}

func NewLogin(
	username valueObject.Username,
	password valueObject.Password,
	ipAddress valueObject.IpAddress,
) Login {
	return Login{
		Username:  username,
		Password:  password,
		IpAddress: ipAddress,
	}
}
