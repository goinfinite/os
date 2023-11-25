package dto

import "github.com/speedianet/os/src/domain/valueObject"

type AddAccount struct {
	Username valueObject.Username `json:"username"`
	Password valueObject.Password `json:"password"`
}

func NewAddAccount(
	username valueObject.Username,
	password valueObject.Password,
) AddAccount {
	return AddAccount{
		Username: username,
		Password: password,
	}
}
