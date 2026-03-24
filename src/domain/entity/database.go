package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type Database struct {
	Name  valueObject.DatabaseName `json:"name"`
	Type  valueObject.DatabaseType `json:"type"`
	Size  tkValueObject.Byte       `json:"size"`
	Users []DatabaseUser           `json:"users"`
}

func NewDatabase(
	name valueObject.DatabaseName,
	dbType valueObject.DatabaseType,
	size tkValueObject.Byte,
	users []DatabaseUser,
) Database {
	return Database{
		Name:  name,
		Type:  dbType,
		Size:  size,
		Users: users,
	}
}
