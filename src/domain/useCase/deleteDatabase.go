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
	_, err := dbQueryRepo.ReadByName(dbName)
	if err != nil {
		return errors.New("DatabaseNotFound")
	}

	err = dbCmdRepo.Delete(dbName)
	if err != nil {
		log.Printf("DeleteDatabaseError: %s", err.Error())
		return errors.New("DeleteDatabaseInfraError")
	}

	log.Printf("Database '%v' deleted.", dbName.String())

	return nil
}
