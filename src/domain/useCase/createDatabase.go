package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func CreateDatabase(
	dbQueryRepo repository.DatabaseQueryRepo,
	dbCmdRepo repository.DatabaseCmdRepo,
	createDatabase dto.CreateDatabase,
) error {
	_, err := dbQueryRepo.ReadByName(createDatabase.DatabaseName)
	if err == nil {
		return errors.New("DatabaseAlreadyExists")
	}

	err = dbCmdRepo.Create(createDatabase.DatabaseName)
	if err != nil {
		slog.Error("CreateDatabaseError", slog.Any("error", err))
		return errors.New("CreateDatabaseInfraError")
	}

	slog.Info(
		"DatabaseCreated",
		slog.String("databaseName", createDatabase.DatabaseName.String()),
	)

	return nil
}
