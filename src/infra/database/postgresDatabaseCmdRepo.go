package databaseInfra

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
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

func (repo PostgresDatabaseCmdRepo) CreateUser(createDatabaseUser dto.CreateDatabaseUser) error {
	dbNameStr := createDatabaseUser.DatabaseName.String()
	dbUserStr := createDatabaseUser.Username.String()

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

	_, err := PostgresqlCmd(
		"GRANT ALL PRIVILEGES ON DATABASE "+dbNameStr+
			" TO "+dbUserStr,
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
		"ALTER DEFAULT PRIVILEGES IN SCHEMA public "+
			"GRANT ALL ON TABLES TO "+dbUserStr,
		&dbNameStr,
	)
	if err != nil {
		return err
	}

	_, err = PostgresqlCmd(
		"ALTER DEFAULT PRIVILEGES IN SCHEMA public "+
			"GRANT ALL ON SEQUENCES TO "+dbUserStr,
		&dbNameStr,
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
		"ALTER DEFAULT PRIVILEGES IN SCHEMA public "+
			"REVOKE ALL ON TABLES FROM "+dbUserStr,
		&dbNameStr,
	)
	if err != nil {
		return err
	}

	_, err = PostgresqlCmd(
		"ALTER DEFAULT PRIVILEGES IN SCHEMA public "+
			"REVOKE ALL ON SEQUENCES FROM "+dbUserStr,
		&dbNameStr,
	)
	if err != nil {
		return err
	}

	postgresDatabaseQueryRepo := PostgresDatabaseQueryRepo{}
	userDbNames, err := postgresDatabaseQueryRepo.GetDatabaseNamesByUser(dbUser)
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
