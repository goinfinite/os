package dto

import "github.com/speedianet/sam/src/domain/valueObject"

type AddDatabaseUser struct {
	DatabaseName valueObject.DatabaseName        `json:"dbName"`
	Username     valueObject.Username            `json:"username"`
	Password     valueObject.Password            `json:"password"`
	Privileges   []valueObject.DatabasePrivilege `json:"privileges"`
}

func NewAddDatabaseUser(
	dbName valueObject.DatabaseName,
	username valueObject.Username,
	password valueObject.Password,
	privileges []valueObject.DatabasePrivilege,
) AddDatabaseUser {
	return AddDatabaseUser{
		DatabaseName: dbName,
		Username:     username,
		Password:     password,
		Privileges:   privileges,
	}
}
