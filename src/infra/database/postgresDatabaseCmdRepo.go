package databaseInfra

import (
	"errors"
	"log"
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
	if err != nil {
		log.Printf("AddDatabaseError: %v", err)
		return errors.New("AddDatabaseError")
	}

	return nil
}

func (repo PostgresDatabaseCmdRepo) Delete(dbName valueObject.DatabaseName) error {
	_, err := PostgresqlCmd(
		"DROP DATABASE " + dbName.String(),
	)
	if err != nil {
		log.Printf("DeleteDatabaseError: %v", err)
		return errors.New("DeleteDatabaseError")
	}

	return nil
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
	if err != nil {
		log.Printf("AddDatabaseUserError: %v", err)
		return errors.New("AddDatabaseUserError")
	}

	return nil
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
	if err != nil {
		log.Printf("DeleteDatabaseUserError: %v", err)
		return errors.New("DeleteDatabaseUserError")
	}

	return nil
}
