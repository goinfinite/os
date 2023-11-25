package databaseInfra

import (
	"github.com/speedianet/os/src/domain/entity"
)

type PostgresDatabaseQueryRepo struct {
}

func (repo PostgresDatabaseQueryRepo) Get() ([]entity.Database, error) {
	return []entity.Database{}, nil
}
