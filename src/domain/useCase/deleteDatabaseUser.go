package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
)

func DeleteDatabaseUser(
	dbQueryRepo repository.DatabaseQueryRepo,
	dbCmdRepo repository.DatabaseCmdRepo,
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteDatabaseUser,
) error {
	_, err := dbQueryRepo.ReadFirst(dto.ReadDatabasesRequest{
		DatabaseName: &deleteDto.DatabaseName,
	})
	if err != nil {
		return errors.New("DatabaseNotFound")
	}

	err = dbCmdRepo.DeleteUser(deleteDto.DatabaseName, deleteDto.Username)
	if err != nil {
		slog.Error("DeleteDatabaseUserError", slog.String("err", err.Error()))
		return errors.New("DeleteDatabaseUserInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteDatabaseUser(deleteDto)

	return nil
}
