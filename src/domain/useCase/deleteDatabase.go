package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
)

func DeleteDatabase(
	dbQueryRepo repository.DatabaseQueryRepo,
	dbCmdRepo repository.DatabaseCmdRepo,
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteDatabase,
) error {
	_, err := dbQueryRepo.ReadFirst(dto.ReadDatabasesRequest{
		DatabaseName: &deleteDto.DatabaseName,
	})
	if err != nil {
		return errors.New("DatabaseNotFound")
	}

	err = dbCmdRepo.Delete(deleteDto.DatabaseName)
	if err != nil {
		slog.Error("DeleteDatabaseError", slog.String("err", err.Error()))
		return errors.New("DeleteDatabaseInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteDatabase(deleteDto)

	return nil
}
