package useCase

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/repository"
)

func GetDatabases(
	databaseQueryRepo repository.DatabaseQueryRepo,
) ([]entity.Database, error) {
	return databaseQueryRepo.Get()
}
