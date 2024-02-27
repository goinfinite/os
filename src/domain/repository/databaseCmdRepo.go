package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
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
