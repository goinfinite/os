package databaseInfra

import (
	"errors"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type PostgresDatabaseCmdRepo struct {
}

func (repo PostgresDatabaseCmdRepo) Add(dbName valueObject.DatabaseName) error {
	return errors.New("NotImplemented")
}

func (repo PostgresDatabaseCmdRepo) Delete(dbName valueObject.DatabaseName) error {
	return errors.New("NotImplemented")
}

func (repo PostgresDatabaseCmdRepo) AddUser(addDatabaseUser dto.AddDatabaseUser) error {
	return errors.New("NotImplemented")
}

func (repo PostgresDatabaseCmdRepo) DeleteUser(
	dbName valueObject.DatabaseName,
	dbUser valueObject.DatabaseUsername,
) error {
	return errors.New("NotImplemented")
}
