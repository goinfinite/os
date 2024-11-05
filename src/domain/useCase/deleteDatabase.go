package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func DeleteDatabase(
	dbQueryRepo repository.DatabaseQueryRepo,
	dbCmdRepo repository.DatabaseCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteDatabase,
) error {
	_, err := dbQueryRepo.ReadByName(deleteDto.DatabaseName)
	if err != nil {
		return errors.New("DatabaseNotFound")
	}

	err = dbCmdRepo.Delete(deleteDto.DatabaseName)
	if err != nil {
		slog.Error("DeleteDatabaseError", slog.Any("error", err))
		return errors.New("DeleteDatabaseInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteDatabase(deleteDto)

	slog.Info(
		"DatabaseDeleted",
		slog.String("databaseName", deleteDto.DatabaseName.String()),
	)

	return nil
}
