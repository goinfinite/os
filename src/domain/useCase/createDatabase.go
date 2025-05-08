package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func CreateDatabase(
	dbQueryRepo repository.DatabaseQueryRepo,
	dbCmdRepo repository.DatabaseCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateDatabase,
) error {
	_, err := dbQueryRepo.ReadFirst(dto.ReadDatabasesRequest{
		DatabaseName: &createDto.DatabaseName,
	})
	if err == nil {
		return errors.New("DatabaseAlreadyExists")
	}

	err = dbCmdRepo.Create(createDto.DatabaseName)
	if err != nil {
		slog.Error("CreateDatabaseError", slog.String("err", err.Error()))
		return errors.New("CreateDatabaseInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateDatabase(createDto)

	return nil
}
