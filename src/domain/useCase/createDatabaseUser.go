package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func CreateDatabaseUser(
	dbQueryRepo repository.DatabaseQueryRepo,
	dbCmdRepo repository.DatabaseCmdRepo,
	createDatabaseUser dto.CreateDatabaseUser,
) error {
	_, err := dbQueryRepo.GetByName(createDatabaseUser.DatabaseName)
	if err != nil {
		return errors.New("DatabaseNotFound")
	}

	if len(createDatabaseUser.Privileges) == 0 {
		createDatabaseUser.Privileges = []valueObject.DatabasePrivilege{
			valueObject.NewDatabasePrivilegePanic("ALL"),
		}
	}

	err = dbCmdRepo.CreateUser(createDatabaseUser)
	if err != nil {
		return errors.New("CreateDatabaseUserError")
	}

	log.Printf(
		"Database user '%s' for '%s' created.",
		createDatabaseUser.Username.String(),
		createDatabaseUser.DatabaseName.String(),
	)

	return nil
}
