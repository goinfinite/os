package useCase

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadDatabases(
	databaseQueryRepo repository.DatabaseQueryRepo,
) ([]entity.Database, error) {
	return databaseQueryRepo.Read()
}
