package databaseInfra

import (
	"errors"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const DatabaseFilePath = "/speedia/sos.db"

type InternalDatabaseService struct {
	Handler *gorm.DB
}

func NewInternalDatabaseService() (*InternalDatabaseService, error) {
	ormSvc, err := gorm.Open(
		sqlite.Open(DatabaseFilePath),
		&gorm.Config{},
	)
	if err != nil {
		return nil, errors.New("DatabaseConnectionError")
	}

	return &InternalDatabaseService{Handler: ormSvc}, nil
}
