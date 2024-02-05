package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type DatabaseCmdRepo interface {
	Add(dbName valueObject.DatabaseName) error
	Delete(dbName valueObject.DatabaseName) error
	AddUser(addDatabaseUser dto.CreateDatabaseUser) error
	DeleteUser(
		dbName valueObject.DatabaseName,
		dbUser valueObject.DatabaseUsername,
	) error
}
