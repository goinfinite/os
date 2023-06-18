package dto

import "github.com/speedianet/sam/src/domain/valueObject"

type Login struct {
	Username valueObject.Username `json:"username"`
	Password valueObject.Password `json:"password"`
}

func NewLogin(
	username valueObject.Username,
	password valueObject.Password,
) Login {
	return Login{
		Username: username,
		Password: password,
	}
}
