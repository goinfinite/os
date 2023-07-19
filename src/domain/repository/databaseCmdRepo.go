package repository

import (
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type DatabaseCmdRepo interface {
	Add(dbName valueObject.DatabaseName) error
	Delete(dbName valueObject.DatabaseName) error
	AddUser(addDatabaseUser dto.AddDatabaseUser) error
}
