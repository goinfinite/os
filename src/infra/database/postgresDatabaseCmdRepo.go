package databaseInfra

import (
	"strings"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type PostgresDatabaseCmdRepo struct {
}

func (repo PostgresDatabaseCmdRepo) Create(dbName valueObject.DatabaseName) error {
	_, err := PostgresqlCmd(
		"CREATE DATABASE " + dbName.String(),
	)

	return err
}

func (repo PostgresDatabaseCmdRepo) Delete(dbName valueObject.DatabaseName) error {
	_, err := PostgresqlCmd(
		"DROP DATABASE " + dbName.String(),
	)

	return err
}

func (repo PostgresDatabaseCmdRepo) CreateUser(createDatabaseUser dto.CreateDatabaseUser) error {
	privilegesStrList := make([]string, len(createDatabaseUser.Privileges))
	for i, privilege := range createDatabaseUser.Privileges {
		privilegesStrList[i] = privilege.String()
	}
	privilegesStr := strings.Join(privilegesStrList, ", ")

	_, err := PostgresqlCmd(
		"GRANT " +
			privilegesStr +
			" ON " +
			createDatabaseUser.DatabaseName.String() +
			".* TO '" +
			createDatabaseUser.Username.String() + "'@'%' " +
			"IDENTIFIED BY '" +
			createDatabaseUser.Password.String() +
			"';",
	)

	return err
}

func (repo PostgresDatabaseCmdRepo) DeleteUser(
	dbName valueObject.DatabaseName,
	dbUser valueObject.DatabaseUsername,
) error {
	_, err := PostgresqlCmd(
		"REVOKE ALL PRIVILEGES ON " +
			dbName.String() +
			".* FROM '" +
			dbUser.String() +
			"'@'%';",
	)

	return err
}
