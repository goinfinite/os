package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func DeleteDatabaseUser(
	dbQueryRepo repository.DatabaseQueryRepo,
	dbCmdRepo repository.DatabaseCmdRepo,
	dbName valueObject.DatabaseName,
	dbUser valueObject.DatabaseUsername,
) error {
	_, err := dbQueryRepo.GetByName(dbName)
	if err != nil {
		return errors.New("DatabaseNotFound")
	}

	err = dbCmdRepo.DeleteUser(dbName, dbUser)
	if err != nil {
		return errors.New("DeleteDatabaseUserError")
	}

	log.Printf(
		"Database user '%s' for '%s' deleted.",
		dbUser.String(),
		dbName.String(),
	)

	return nil
}
