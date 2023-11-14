package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
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

	if len(addDatabaseUser.Privileges) == 0 {
		addDatabaseUser.Privileges = []valueObject.DatabasePrivilege{
			valueObject.NewDatabasePrivilegePanic("ALL"),
		}
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
