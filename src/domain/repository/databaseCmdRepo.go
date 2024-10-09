package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type DatabaseCmdRepo interface {
	Create(dbName valueObject.DatabaseName) error
	Delete(dbName valueObject.DatabaseName) error
	CreateUser(createDatabaseUser dto.CreateDatabaseUser) error
	DeleteUser(
		dbName valueObject.DatabaseName,
		dbUser valueObject.DatabaseUsername,
	) error
}
