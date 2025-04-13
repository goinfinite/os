package databaseInfra

import (
	"errors"
	"regexp"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type PostgresDatabaseCmdRepo struct {
}

func (repo PostgresDatabaseCmdRepo) Create(dbName valueObject.DatabaseName) error {
	_, err := PostgresqlCmd(
		"CREATE DATABASE "+dbName.String(),
		nil,
	)

	return err
}

func (repo PostgresDatabaseCmdRepo) Delete(dbName valueObject.DatabaseName) error {
	_, err := PostgresqlCmd(
		"DROP DATABASE "+dbName.String(),
		nil,
	)

	return err
}

func (repo PostgresDatabaseCmdRepo) CreateUser(
	createDatabaseUser dto.CreateDatabaseUser,
) error {
	dbUserStr := createDatabaseUser.Username.String()
	if regexp.MustCompile(`^\d`).MatchString(dbUserStr) {
		return errors.New("PostgresUsernameCannotStartWithNumbers")
	}

	postgresDatabaseQueryRepo := PostgresDatabaseQueryRepo{}
	userExists := postgresDatabaseQueryRepo.UserExists(createDatabaseUser.Username)
	if !userExists {
		_, err := PostgresqlCmd(
			"CREATE USER "+dbUserStr+" WITH PASSWORD '"+createDatabaseUser.Password.String()+"'",
			nil,
		)
		if err != nil {
			return err
		}
	}

	dbNameStr := createDatabaseUser.DatabaseName.String()
	_, err := PostgresqlCmd(
		"GRANT ALL PRIVILEGES ON DATABASE "+dbNameStr+" TO "+dbUserStr,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = PostgresqlCmd("GRANT ALL ON ALL TABLES IN SCHEMA public TO "+dbUserStr, nil)
	if err != nil {
		return err
	}

	_, err = PostgresqlCmd("GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO "+dbUserStr, nil)
	if err != nil {
		return err
	}

	_, err = PostgresqlCmd(
		"ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO "+dbUserStr,
		&createDatabaseUser.DatabaseName,
	)
	if err != nil {
		return err
	}

	_, err = PostgresqlCmd(
		"ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO "+dbUserStr,
		&createDatabaseUser.DatabaseName,
	)
	return err
}

func (repo PostgresDatabaseCmdRepo) DeleteUser(
	dbName valueObject.DatabaseName,
	dbUser valueObject.DatabaseUsername,
) error {
	dbNameStr := dbName.String()
	dbUserStr := dbUser.String()

	_, err := PostgresqlCmd(
		"REVOKE ALL ON DATABASE "+dbNameStr+" FROM "+dbUserStr,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = PostgresqlCmd(
		"ALTER DEFAULT PRIVILEGES IN SCHEMA public REVOKE ALL ON TABLES FROM "+dbUserStr,
		&dbName,
	)
	if err != nil {
		return err
	}

	_, err = PostgresqlCmd(
		"ALTER DEFAULT PRIVILEGES IN SCHEMA public REVOKE ALL ON SEQUENCES FROM "+dbUserStr,
		&dbName,
	)
	if err != nil {
		return err
	}

	postgresDatabaseQueryRepo := PostgresDatabaseQueryRepo{}
	userDbNames, err := postgresDatabaseQueryRepo.ReadDatabaseNamesByUser(dbUser)
	if err != nil {
		return err
	}

	userStillHasPermissions := len(userDbNames) > 0
	if userStillHasPermissions {
		return nil
	}

	_, err = PostgresqlCmd("DROP USER "+dbUserStr, nil)
	return err
}
