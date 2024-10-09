package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func CreateDatabaseUser(
	dbQueryRepo repository.DatabaseQueryRepo,
	dbCmdRepo repository.DatabaseCmdRepo,
	createDatabaseUser dto.CreateDatabaseUser,
) error {
	_, err := dbQueryRepo.ReadByName(createDatabaseUser.DatabaseName)
	if err != nil {
		return errors.New("DatabaseNotFound")
	}

	if len(createDatabaseUser.Privileges) == 0 {
		defaultPrivilege, err := valueObject.NewDatabasePrivilege("ALL")
		if err != nil {
			return err
		}

		createDatabaseUser.Privileges = []valueObject.DatabasePrivilege{
			defaultPrivilege,
		}
	}

	err = dbCmdRepo.CreateUser(createDatabaseUser)
	if err != nil {
		slog.Error("CreateDatabaseUserError", slog.Any("error", err))
		return errors.New("CreateDatabaseUserInfraError")
	}

	slog.Info(
		"DatabaseUserCreated",
		slog.String("databaseName", createDatabaseUser.DatabaseName.String()),
		slog.String("databaseUsername", createDatabaseUser.Username.String()),
	)

	return nil
}
