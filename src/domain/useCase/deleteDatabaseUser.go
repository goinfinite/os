package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func DeleteDatabaseUser(
	dbQueryRepo repository.DatabaseQueryRepo,
	dbCmdRepo repository.DatabaseCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteDatabaseUser,
) error {
	_, err := dbQueryRepo.ReadByName(deleteDto.DatabaseName)
	if err != nil {
		return errors.New("DatabaseNotFound")
	}

	err = dbCmdRepo.DeleteUser(deleteDto.DatabaseName, deleteDto.Username)
	if err != nil {
		slog.Error("DeleteDatabaseUserError", slog.Any("error", err))
		return errors.New("DeleteDatabaseUserInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteDatabaseUser(deleteDto)

	slog.Info(
		"DatabaseUserDeleted",
		slog.String("databaseName", deleteDto.DatabaseName.String()),
		slog.String("databaseUsername", deleteDto.Username.String()),
	)

	return nil
}
