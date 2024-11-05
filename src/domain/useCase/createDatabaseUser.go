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
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateDatabaseUser,
) error {
	_, err := dbQueryRepo.ReadByName(createDto.DatabaseName)
	if err != nil {
		return errors.New("DatabaseNotFound")
	}

	if len(createDto.Privileges) == 0 {
		defaultPrivilege, err := valueObject.NewDatabasePrivilege("ALL")
		if err != nil {
			return err
		}

		createDto.Privileges = []valueObject.DatabasePrivilege{
			defaultPrivilege,
		}
	}

	err = dbCmdRepo.CreateUser(createDto)
	if err != nil {
		slog.Error("CreateDatabaseUserError", slog.Any("error", err))
		return errors.New("CreateDatabaseUserInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateDatabaseUser(createDto)

	slog.Info(
		"DatabaseUserCreated",
		slog.String("databaseName", createDto.DatabaseName.String()),
		slog.String("databaseUsername", createDto.Username.String()),
	)

	return nil
}
