package databaseInfra

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type PostgresDatabaseCmdRepo struct {
}

func (repo PostgresDatabaseCmdRepo) Add(dbName valueObject.DatabaseName) error {
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

func (repo PostgresDatabaseCmdRepo) addPermissionsToUser(
	dbName valueObject.DatabaseName,
	dbUser valueObject.DatabaseUsername,
) error {
	_, err := PostgresqlCmd(
		"GRANT ALL PRIVILEGES ON DATABASE "+dbName.String()+
			" TO "+dbUser.String(),
		nil,
	)
	if err != nil {
		return err
	}

	_, err = PostgresqlCmd("GRANT ALL ON ALL TABLES IN SCHEMA public TO "+dbUser.String(), nil)
	if err != nil {
		return err
	}

	_, err = PostgresqlCmd("GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO "+dbUser.String(), nil)
	return err
}

func (repo PostgresDatabaseCmdRepo) setUserDefaultPermissions(
	dbName valueObject.DatabaseName,
	dbUser valueObject.DatabaseUsername,
) error {
	dbNameStr := dbName.String()

	_, err := PostgresqlCmd(
		"ALTER DEFAULT PRIVILEGES IN SCHEMA public "+
			"GRANT ALL ON TABLES TO "+dbUser.String(),
		&dbNameStr,
	)
	if err != nil {
		return err
	}

	_, err = PostgresqlCmd(
		"ALTER DEFAULT PRIVILEGES IN SCHEMA public "+
			"GRANT ALL ON SEQUENCES TO "+dbUser.String(),
		&dbNameStr,
	)

	return err
}

func (repo PostgresDatabaseCmdRepo) AddUser(addDatabaseUser dto.AddDatabaseUser) error {
	_, err := PostgresqlCmd(
		"CREATE USER "+addDatabaseUser.Username.String()+
			" WITH PASSWORD '"+addDatabaseUser.Password.String()+"'",
		nil,
	)
	if err != nil {
		return err
	}

	err = repo.addPermissionsToUser(addDatabaseUser.DatabaseName, addDatabaseUser.Username)
	if err != nil {
		return err
	}

	err = repo.setUserDefaultPermissions(addDatabaseUser.DatabaseName, addDatabaseUser.Username)
	return err
}

func (repo PostgresDatabaseCmdRepo) revokeUserDefaultPermissions(
	dbName valueObject.DatabaseName,
	dbUser valueObject.DatabaseUsername,
) error {
	dbNameStr := dbName.String()

	_, err := PostgresqlCmd(
		"ALTER DEFAULT PRIVILEGES IN SCHEMA public "+
			"REVOKE ALL ON TABLES FROM "+dbUser.String(),
		&dbNameStr,
	)
	if err != nil {
		return err
	}

	_, err = PostgresqlCmd(
		"ALTER DEFAULT PRIVILEGES IN SCHEMA public "+
			"REVOKE ALL ON SEQUENCES FROM "+dbUser.String(),
		&dbNameStr,
	)
	return err
}

func (repo PostgresDatabaseCmdRepo) DeleteUser(
	dbName valueObject.DatabaseName,
	dbUser valueObject.DatabaseUsername,
) error {
	_, err := PostgresqlCmd(
		"REVOKE ALL ON DATABASE "+dbName.String()+" FROM "+dbUser.String(),
		nil,
	)
	if err != nil {
		return err
	}

	err = repo.revokeUserDefaultPermissions(dbName, dbUser)
	if err != nil {
		return err
	}

	_, err = PostgresqlCmd("DROP USER "+dbUser.String(), nil)
	return err
}
