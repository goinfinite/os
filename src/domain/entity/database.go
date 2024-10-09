package entity

import "github.com/goinfinite/os/src/domain/valueObject"

type Database struct {
	Name  valueObject.DatabaseName `json:"name"`
	Type  valueObject.DatabaseType `json:"type"`
	Size  valueObject.Byte         `json:"size"`
	Users []DatabaseUser           `json:"users"`
}

func NewDatabase(
	name valueObject.DatabaseName,
	dbType valueObject.DatabaseType,
	size valueObject.Byte,
	users []DatabaseUser,
) Database {
	return Database{
		Name:  name,
		Type:  dbType,
		Size:  size,
		Users: users,
	}
}
