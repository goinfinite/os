package internalDbInfra

import (
	"errors"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const DatabaseFilePath = "/speedia/sos.db"

type PersistentDatabaseService struct {
	Handler *gorm.DB
}

func NewPersistentDatabaseService() (*PersistentDatabaseService, error) {
	ormSvc, err := gorm.Open(
		sqlite.Open(DatabaseFilePath),
		&gorm.Config{},
	)
	if err != nil {
		return nil, errors.New("DatabaseConnectionError")
	}

	return &PersistentDatabaseService{Handler: ormSvc}, nil
}
