package databaseInfra

import (
	"errors"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
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

func (repo DatabaseQueryRepo) Get() ([]entity.Database, error) {
	switch repo.dbType {
	case "mariadb":
		return MysqlDatabaseQueryRepo{}.Get()
	case "postgresql":
		return PostgresDatabaseQueryRepo{}.Get()
	default:
		return []entity.Database{}, errors.New("DatabaseTypeNotSupported")
	}
}

func (repo DatabaseQueryRepo) GetByName(
	name valueObject.DatabaseName,
) (entity.Database, error) {
	dbs, err := repo.Get()
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
