package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadDatabases(
	databaseQueryRepo repository.DatabaseQueryRepo,
) (databasesList []entity.Database, err error) {
	databasesList, err = databaseQueryRepo.Read()
	if err != nil {
		slog.Error("ReadDatabasesError", slog.Any("err", err))
		return databasesList, errors.New("ReadDatabasesInfraError")
	}

	return databasesList, nil
}
