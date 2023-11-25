package useCase

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func GetDatabases(
	databaseQueryRepo repository.DatabaseQueryRepo,
) ([]entity.Database, error) {
	return databaseQueryRepo.Get()
}
