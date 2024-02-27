package databaseInfra

import (
	"errors"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type DatabaseCmdRepo struct {
	dbType valueObject.DatabaseType
}

func NewDatabaseCmdRepo(
	dbType valueObject.DatabaseType,
) *DatabaseCmdRepo {
	return &DatabaseCmdRepo{
		dbType: dbType,
	}
}

func (repo DatabaseCmdRepo) Create(dbName valueObject.DatabaseName) error {
	switch repo.dbType {
	case "mysql":
		return MysqlDatabaseCmdRepo{}.Create(dbName)
	case "postgres":
		return PostgresDatabaseCmdRepo{}.Create(dbName)
	default:
		return errors.New("DatabaseTypeNotSupported")
	}
}

func (repo DatabaseCmdRepo) Delete(dbName valueObject.DatabaseName) error {
	switch repo.dbType {
	case "mysql":
		return MysqlDatabaseCmdRepo{}.Delete(dbName)
	case "postgres":
		return PostgresDatabaseCmdRepo{}.Delete(dbName)
	default:
		return errors.New("DatabaseTypeNotSupported")
	}
}

func (repo DatabaseCmdRepo) CreateUser(createDatabaseUser dto.CreateDatabaseUser) error {
	switch repo.dbType {
	case "mysql":
		return MysqlDatabaseCmdRepo{}.CreateUser(createDatabaseUser)
	case "postgres":
		return PostgresDatabaseCmdRepo{}.CreateUser(createDatabaseUser)
	default:
		return errors.New("DatabaseTypeNotSupported")
	}
}

func (repo DatabaseCmdRepo) DeleteUser(
	dbName valueObject.DatabaseName,
	dbUser valueObject.DatabaseUsername,
) error {
	switch repo.dbType {
	case "mysql":
		return MysqlDatabaseCmdRepo{}.DeleteUser(dbName, dbUser)
	case "postgres":
		return PostgresDatabaseCmdRepo{}.DeleteUser(dbName, dbUser)
	default:
		return errors.New("DatabaseTypeNotSupported")
	}
}
