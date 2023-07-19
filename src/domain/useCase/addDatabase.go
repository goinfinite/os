package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
)

func AddDatabase(
	dbQueryRepo repository.DatabaseQueryRepo,
	dbCmdRepo repository.DatabaseCmdRepo,
	addDatabase dto.AddDatabase,
) error {
	_, err := dbQueryRepo.GetByName(addDatabase.DatabaseName)
	if err == nil {
		return errors.New("DatabaseAlreadyExists")
	}

	err = dbCmdRepo.Add(addDatabase.DatabaseName)
	if err != nil {
		return errors.New("AddDatabaseError")
	}

	log.Printf("Database %s added", addDatabase.DatabaseName)

	return nil
}
