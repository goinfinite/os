package databaseInfra

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type MysqlDatabaseCmdRepo struct {
}

func (repo MysqlDatabaseCmdRepo) Create(dbName valueObject.DatabaseName) error {
	_, err := MysqlCmd(
		"CREATE DATABASE " + dbName.String(),
	)
	if err != nil {
		slog.Error("CreateDatabaseError", slog.String("err", err.Error()))
		return errors.New("CreateDatabaseError")
	}

	return nil
}

func (repo MysqlDatabaseCmdRepo) Delete(dbName valueObject.DatabaseName) error {
	_, err := MysqlCmd(
		"DROP DATABASE " + dbName.String(),
	)
	if err != nil {
		slog.Error("DeleteDatabaseError", slog.String("err", err.Error()))
		return errors.New("DeleteDatabaseError")
	}

	return nil
}

func (repo MysqlDatabaseCmdRepo) CreateUser(createDatabaseUser dto.CreateDatabaseUser) error {
	privilegesStrList := make([]string, len(createDatabaseUser.Privileges))
	for i, privilege := range createDatabaseUser.Privileges {
		privilegesStrList[i] = privilege.String()
	}
	privilegesStr := strings.Join(privilegesStrList, ", ")

	_, err := MysqlCmd(
		"GRANT " +
			privilegesStr +
			" ON " +
			createDatabaseUser.DatabaseName.String() +
			".* TO '" +
			createDatabaseUser.Username.String() + "'@'%' " +
			"IDENTIFIED BY '" +
			createDatabaseUser.Password.String() +
			"'; " +
			"FLUSH PRIVILEGES;",
	)
	if err != nil {
		slog.Error("CreateDatabaseUserError", slog.String("err", err.Error()))
		return errors.New("CreateDatabaseUserError")
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
		slog.Error("DeleteDatabaseUserError", slog.String("err", err.Error()))
		return errors.New("DeleteDatabaseUserError")
	}

	return nil
}
