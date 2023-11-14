package entity

import "github.com/speedianet/os/src/domain/valueObject"

type DatabaseUser struct {
	Username   valueObject.DatabaseUsername    `json:"username"`
	DbName     valueObject.DatabaseName        `json:"dbName"`
	DbType     valueObject.DatabaseType        `json:"dbType"`
	Privileges []valueObject.DatabasePrivilege `json:"privileges"`
}

func NewDatabaseUser(
	username valueObject.DatabaseUsername,
	dbName valueObject.DatabaseName,
	dbType valueObject.DatabaseType,
	privileges []valueObject.DatabasePrivilege,
) DatabaseUser {
	return DatabaseUser{
		Username:   username,
		DbName:     dbName,
		DbType:     dbType,
		Privileges: privileges,
	}
}
