package useCase

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func ReadDatabases(
	databaseQueryRepo repository.DatabaseQueryRepo,
) ([]entity.Database, error) {
	return databaseQueryRepo.Read()
}
