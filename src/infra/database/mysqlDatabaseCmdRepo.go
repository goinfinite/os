package databaseInfra

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/valueObject"
)

type MysqlDatabaseCmdRepo struct {
}

func (repo MysqlDatabaseCmdRepo) Add(dbName valueObject.DatabaseName) error {
	_, err := MysqlCmd(
		"CREATE DATABASE " + dbName.String(),
	)
	if err != nil {
		log.Printf("AddDatabaseError: %v", err)
		return errors.New("AddDatabaseError")
	}

	return nil
}
