package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func CreateDatabase(
	dbQueryRepo repository.DatabaseQueryRepo,
	dbCmdRepo repository.DatabaseCmdRepo,
	createDatabase dto.CreateDatabase,
) error {
	_, err := dbQueryRepo.GetByName(createDatabase.DatabaseName)
	if err == nil {
		return errors.New("DatabaseAlreadyExists")
	}

	err = dbCmdRepo.Create(createDatabase.DatabaseName)
	if err != nil {
		log.Printf("CreateDatabaseError: %s", err.Error())
		return errors.New("CreateDatabaseInfraError")
	}

	log.Printf("Database %s created", createDatabase.DatabaseName)

	return nil
}
