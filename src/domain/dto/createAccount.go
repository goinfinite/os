package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateAccount struct {
	Username valueObject.Username `json:"username"`
	Password valueObject.Password `json:"password"`
}

func NewCreateAccount(
	username valueObject.Username,
	password valueObject.Password,
) CreateAccount {
	return CreateAccount{
		Username: username,
		Password: password,
	}
}
