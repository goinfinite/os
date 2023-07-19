package infra

import (
	"errors"

	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
	databaseInfra "github.com/speedianet/sam/src/infra/database"
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
	case "mysql":
		return databaseInfra.MysqlDatabaseQueryRepo{}.Get()
	case "postgres":
		return databaseInfra.PostgresDatabaseQueryRepo{}.Get()
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
