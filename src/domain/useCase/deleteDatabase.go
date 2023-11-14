package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func DeleteDatabase(
	dbQueryRepo repository.DatabaseQueryRepo,
	dbCmdRepo repository.DatabaseCmdRepo,
	dbName valueObject.DatabaseName,
) error {
	_, err := dbQueryRepo.GetByName(dbName)
	if err != nil {
		return errors.New("DatabaseNotFound")
	}

	err = dbCmdRepo.Delete(dbName)
	if err != nil {
		return errors.New("DeleteDatabaseError")
	}

	log.Printf("Database '%v' deleted.", dbName.String())

	return nil
}
