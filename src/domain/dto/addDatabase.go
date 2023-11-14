package dto

import "github.com/speedianet/os/src/domain/valueObject"

type AddDatabase struct {
	DatabaseName valueObject.DatabaseName `json:"dbName"`
}

func NewAddDatabase(
	dbName valueObject.DatabaseName,
) AddDatabase {
	return AddDatabase{
		DatabaseName: dbName,
	}
}
