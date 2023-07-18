package infra

import (
	"github.com/speedianet/sam/src/domain/repository"
	"github.com/speedianet/sam/src/domain/valueObject"
	databaseInfra "github.com/speedianet/sam/src/infra/database"
)

type DatabaseQueryRepo struct {
}

func NewDatabaseQueryRepo(
	dbType valueObject.DatabaseType,
) repository.DatabaseQueryRepo {
	if dbType.String() == "postgres" {
		return databaseInfra.PostgresDatabaseQueryRepo{}
	}
	return databaseInfra.MysqlDatabaseQueryRepo{}
}
