package databaseInfra

import (
	"errors"
	"log"
	"strings"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type MysqlDatabaseCmdRepo struct {
}

func (repo MysqlDatabaseCmdRepo) Add(dbName valueObject.DatabaseName) error {
	_, err := MysqlCmd(
		"CREATE DATABASE " + dbName.String(),
	)
	if err != nil {
		log.Printf("AddDatabaseError: %v", err)
		return errors.New("AddDatabaseError")
	}

	return nil
}

func (repo MysqlDatabaseCmdRepo) Delete(dbName valueObject.DatabaseName) error {
	_, err := MysqlCmd(
		"DROP DATABASE " + dbName.String(),
	)
	if err != nil {
		log.Printf("DeleteDatabaseError: %v", err)
		return errors.New("DeleteDatabaseError")
	}

	return nil
}

func (repo MysqlDatabaseCmdRepo) AddUser(addDatabaseUser dto.AddDatabaseUser) error {
	privilegesStrList := make([]string, len(addDatabaseUser.Privileges))
	for i, privilege := range addDatabaseUser.Privileges {
		privilegesStrList[i] = privilege.String()
	}
	privilegesStr := strings.Join(privilegesStrList, ", ")

	_, err := MysqlCmd(
		"GRANT " +
			privilegesStr +
			" PRIVILEGES ON " +
			addDatabaseUser.DatabaseName.String() +
			".* TO '" +
			addDatabaseUser.Username.String() + "'@'%' " +
			"IDENTIFIED BY '" +
			addDatabaseUser.Password.String() +
			"'; " +
			"FLUSH PRIVILEGES;",
	)
	if err != nil {
		log.Printf("AddDatabaseUserError: %v", err)
		return errors.New("AddDatabaseUserError")
	}

	return nil
}

func (repo MysqlDatabaseCmdRepo) DeleteUser(
	dbName valueObject.DatabaseName,
	dbUser valueObject.DatabaseUsername,
) error {
	_, err := MysqlCmd(
		"REVOKE ALL PRIVILEGES ON " +
			dbName.String() +
			".* FROM '" +
			dbUser.String() +
			"'@'%'; " +
			"FLUSH PRIVILEGES;",
	)
	if err != nil {
		log.Printf("DeleteDatabaseUserError: %v", err)
		return errors.New("DeleteDatabaseUserError")
	}

	return nil
}
