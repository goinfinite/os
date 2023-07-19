package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
)

func AddDatabaseUser(
	dbQueryRepo repository.DatabaseQueryRepo,
	dbCmdRepo repository.DatabaseCmdRepo,
	addDatabaseUser dto.AddDatabaseUser,
) error {
	_, err := dbQueryRepo.GetByName(addDatabaseUser.DatabaseName)
	if err != nil {
		return errors.New("DatabaseNotFound")
	}

	err = dbCmdRepo.AddUser(addDatabaseUser)
	if err != nil {
		return errors.New("AddDatabaseUserError")
	}

	log.Printf(
		"Database user '%s' for '%s' added.",
		addDatabaseUser.Username.String(),
		addDatabaseUser.DatabaseName.String(),
	)

	return nil
}
