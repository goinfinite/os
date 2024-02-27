package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CreateDatabase struct {
	DatabaseName valueObject.DatabaseName `json:"dbName"`
}

func NewCreateDatabase(
	dbName valueObject.DatabaseName,
) CreateDatabase {
	return CreateDatabase{
		DatabaseName: dbName,
	}
}
