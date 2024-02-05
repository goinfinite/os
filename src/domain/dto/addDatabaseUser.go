package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CreateDatabaseUser struct {
	DatabaseName valueObject.DatabaseName        `json:"dbName"`
	Username     valueObject.DatabaseUsername    `json:"username"`
	Password     valueObject.Password            `json:"password"`
	Privileges   []valueObject.DatabasePrivilege `json:"privileges"`
}

func NewCreateDatabaseUser(
	dbName valueObject.DatabaseName,
	username valueObject.DatabaseUsername,
	password valueObject.Password,
	privileges []valueObject.DatabasePrivilege,
) CreateDatabaseUser {
	return CreateDatabaseUser{
		DatabaseName: dbName,
		Username:     username,
		Password:     password,
		Privileges:   privileges,
	}
}
