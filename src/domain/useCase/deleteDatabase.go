package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func DeleteDatabase(
	dbQueryRepo repository.DatabaseQueryRepo,
	dbCmdRepo repository.DatabaseCmdRepo,
	dbName valueObject.DatabaseName,
) error {
	_, err := dbQueryRepo.ReadByName(dbName)
	if err != nil {
		return errors.New("DatabaseNotFound")
	}

	err = dbCmdRepo.Delete(dbName)
	if err != nil {
		slog.Error("DeleteDatabaseError", slog.Any("error", err))
		return errors.New("DeleteDatabaseInfraError")
	}

	slog.Info(
		"DatabaseDeleted",
		slog.String("databaseName", dbName.String()),
	)

	return nil
}
