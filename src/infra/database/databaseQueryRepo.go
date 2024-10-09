package databaseInfra

import (
	"errors"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type DatabaseQueryRepo struct {
	dbType valueObject.DatabaseType
}

func NewDatabaseQueryRepo(
	dbType valueObject.DatabaseType,
) *DatabaseQueryRepo {
	return &DatabaseQueryRepo{
		dbType: dbType,
	}
}

func (repo DatabaseQueryRepo) Read() ([]entity.Database, error) {
	switch repo.dbType {
	case "mariadb":
		return MysqlDatabaseQueryRepo{}.Read()
	case "postgresql":
		return PostgresDatabaseQueryRepo{}.Read()
	default:
		return []entity.Database{}, errors.New("DatabaseTypeNotSupported")
	}
}

func (repo DatabaseQueryRepo) ReadByName(
	name valueObject.DatabaseName,
) (entity.Database, error) {
	dbs, err := repo.Read()
	if err != nil {
		return entity.Database{}, err
	}

	for _, db := range dbs {
		if db.Name == name {
			return db, nil
		}
	}

	return entity.Database{}, errors.New("DatabaseNotFound")
}
