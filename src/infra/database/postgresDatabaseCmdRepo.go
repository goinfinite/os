package databaseInfra

import (
	"strings"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type PostgresDatabaseCmdRepo struct {
}

func (repo PostgresDatabaseCmdRepo) Add(dbName valueObject.DatabaseName) error {
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

func (repo PostgresDatabaseCmdRepo) AddUser(addDatabaseUser dto.AddDatabaseUser) error {
	privilegesStrList := make([]string, len(addDatabaseUser.Privileges))
	for i, privilege := range addDatabaseUser.Privileges {
		privilegesStrList[i] = privilege.String()
	}
	privilegesStr := strings.Join(privilegesStrList, ", ")

	_, err := PostgresqlCmd(
		"GRANT " +
			privilegesStr +
			" ON " +
			addDatabaseUser.DatabaseName.String() +
			".* TO '" +
			addDatabaseUser.Username.String() + "'@'%' " +
			"IDENTIFIED BY '" +
			addDatabaseUser.Password.String() +
			"';",
	)

	return err
}

func (repo PostgresDatabaseCmdRepo) DeleteUser(
	dbName valueObject.DatabaseName,
	dbUser valueObject.DatabaseUsername,
) error {
	_, err := MysqlCmd(
		"REVOKE ALL PRIVILEGES ON " +
			dbName.String() +
			".* FROM '" +
			dbUser.String() +
			"'@'%';",
	)

	return err
}
