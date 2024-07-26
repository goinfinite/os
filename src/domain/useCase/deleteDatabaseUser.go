package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func DeleteDatabaseUser(
	dbQueryRepo repository.DatabaseQueryRepo,
	dbCmdRepo repository.DatabaseCmdRepo,
	dbName valueObject.DatabaseName,
	dbUser valueObject.DatabaseUsername,
) error {
	_, err := dbQueryRepo.ReadByName(dbName)
	if err != nil {
		return errors.New("DatabaseNotFound")
	}

	err = dbCmdRepo.DeleteUser(dbName, dbUser)
	if err != nil {
		slog.Error("DeleteDatabaseUserError", slog.Any("error", err))
		return errors.New("DeleteDatabaseUserInfraError")
	}

	slog.Info(
		"DatabaseUserDeleted",
		slog.String("databaseName", dbUser.String()),
		slog.String("databaseUsername", dbName.String()),
	)

	return nil
}
