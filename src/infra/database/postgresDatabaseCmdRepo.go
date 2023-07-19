package databaseInfra

import (
	"errors"

	"github.com/speedianet/sam/src/domain/valueObject"
)

type PostgresDatabaseCmdRepo struct {
}

func (repo PostgresDatabaseCmdRepo) Add(dbName valueObject.DatabaseName) error {
	return errors.New("NotImplemented")
}

func (repo PostgresDatabaseCmdRepo) Delete(dbName valueObject.DatabaseName) error {
	return errors.New("NotImplemented")
}
