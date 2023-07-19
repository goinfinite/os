package dto

import "github.com/speedianet/sam/src/domain/valueObject"

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
