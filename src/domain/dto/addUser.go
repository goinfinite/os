package dto

import "github.com/speedianet/sam/src/domain/valueObject"

type AddUser struct {
	Username valueObject.Username `json:"username"`
	Password valueObject.Password `json:"password"`
}

func NewAddUser(
	username valueObject.Username,
	password valueObject.Password,
) AddUser {
	return AddUser{
		Username: username,
		Password: password,
	}
}
