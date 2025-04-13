package databaseInfra

import (
	"errors"
	"log/slog"
	"regexp"
	"strings"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
)

type PostgresDatabaseQueryRepo struct {
}

func PostgresqlCmd(cmd string, dbName *valueObject.DatabaseName) (string, error) {
	cmdArgs := []string{"-U", "postgres", "-tAc", cmd}
	if dbName != nil {
		cmdArgs = append(cmdArgs, "-d", dbName.String())
	}

	return infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "psql",
		Args:    cmdArgs,
	})
}

func (repo PostgresDatabaseQueryRepo) readDatabaseNames() ([]valueObject.DatabaseName, error) {
	databaseNames := []valueObject.DatabaseName{}

	rawDatabaseNames, err := PostgresqlCmd(
		"SELECT datname FROM pg_database WHERE datistemplate = false AND datname != 'postgres'",
		nil,
	)
	if err != nil {
		return databaseNames, errors.New("ReadDatabaseNamesError: " + err.Error())
	}

	for rawDatabaseName := range strings.SplitSeq(rawDatabaseNames, "\n") {
		rawDatabaseName = strings.TrimSpace(rawDatabaseName)
		if rawDatabaseName == "" {
			continue
		}

		dbName, err := valueObject.NewDatabaseName(rawDatabaseName)
		if err != nil {
			slog.Debug(
				err.Error(),
				slog.String("rawDbName", rawDatabaseName),
			)
			continue
		}
		databaseNames = append(databaseNames, dbName)
	}

	return databaseNames, nil
}

func (repo PostgresDatabaseQueryRepo) readDatabaseSize(
	dbName valueObject.DatabaseName,
) (valueObject.Byte, error) {
	rawDatabaseSize, err := PostgresqlCmd(
		"SELECT pg_database_size('"+dbName.String()+"')",
		nil,
	)
	if err != nil {
		return 0, errors.New("ReadDatabaseSizeError: " + err.Error())
	}

	return valueObject.NewByte(rawDatabaseSize)
}

func (repo PostgresDatabaseQueryRepo) readDatabaseUsernames(
	dbName valueObject.DatabaseName,
) ([]valueObject.DatabaseUsername, error) {
	dbUsernames := []valueObject.DatabaseUsername{}

	rawDatabaseUsers, err := PostgresqlCmd(
		"SELECT datacl FROM pg_database WHERE datname = '"+dbName.String()+"'",
		nil,
	)
	if err != nil {
		return dbUsernames, errors.New("ReadDatabaseUsersError: " + err.Error())
	}

	dbUsersRegex := regexp.MustCompile(`(?P<username>[\w\-]{2,256})(?:=)`)
	rawDbUsersMatches := dbUsersRegex.FindAllStringSubmatch(rawDatabaseUsers, -1)
	if len(rawDbUsersMatches) == 0 {
		return dbUsernames, nil
	}

	defaultDbUser := "postgres"
	for _, rawDbUserMatch := range rawDbUsersMatches {
		if len(rawDbUserMatch) < 2 {
			continue
		}

		rawDatabaseUsername := rawDbUserMatch[1]
		if rawDatabaseUsername == defaultDbUser {
			continue
		}

		dbUsername, err := valueObject.NewDatabaseUsername(rawDatabaseUsername)
		if err != nil {
			slog.Debug(
				err.Error(),
				slog.String("rawDbUser", rawDatabaseUsername),
			)
			continue
		}

		dbUsernames = append(dbUsernames, dbUsername)
	}

	return dbUsernames, nil
}

func (repo PostgresDatabaseQueryRepo) UserExists(dbUsername valueObject.DatabaseUsername) bool {
	rawDatabaseUsers, err := PostgresqlCmd(
		"SELECT rolname FROM pg_roles WHERE rolname = '"+dbUsername.String()+"'",
		nil,
	)
	if err != nil {
		return false
	}

	return strings.TrimSpace(rawDatabaseUsers) != ""
}

func (repo PostgresDatabaseQueryRepo) ReadDatabaseNamesByUser(
	dbUsername valueObject.DatabaseUsername,
) ([]valueObject.DatabaseName, error) {
	databaseNames := []valueObject.DatabaseName{}

	rawDatabaseNames, err := PostgresqlCmd(
		"SELECT d.datname FROM pg_database d JOIN pg_roles r ON r.oid = d.datdba WHERE r.rolname = '"+dbUsername.String()+"' UNION SELECT d.datname FROM pg_database d WHERE EXISTS (SELECT 1 FROM unnest(d.datacl::text[]) acl WHERE acl LIKE '%"+dbUsername.String()+"%')",
		nil,
	)
	if err != nil {
		return databaseNames, errors.New("ReadDatabaseNamesByUserError: " + err.Error())
	}

	for rawDatabaseName := range strings.SplitSeq(rawDatabaseNames, "\n") {
		rawDatabaseName = strings.TrimSpace(rawDatabaseName)
		if rawDatabaseName == "" {
			continue
		}

		dbName, err := valueObject.NewDatabaseName(rawDatabaseName)
		if err != nil {
			slog.Debug(
				err.Error(),
				slog.String("rawDbName", rawDatabaseName),
			)
			continue
		}

		databaseNames = append(databaseNames, dbName)
	}

	return databaseNames, nil
}

func (repo PostgresDatabaseQueryRepo) readAllDatabases() ([]entity.Database, error) {
	databaseEntities := []entity.Database{}

	databaseNames, err := repo.readDatabaseNames()
	if err != nil {
		return databaseEntities, errors.New("ReadDatabaseNamesError: " + err.Error())
	}
	dbType, _ := valueObject.NewDatabaseType("postgresql")

	for _, dbName := range databaseNames {
		dbSize, err := repo.readDatabaseSize(dbName)
		if err != nil {
			dbSize, _ = valueObject.NewByte(0)
		}

		dbUsernames, err := repo.readDatabaseUsernames(dbName)
		if err != nil {
			slog.Debug(
				"ReadDatabaseUsersError",
				slog.String("dbName", dbName.String()),
				slog.String("err", err.Error()),
			)
			continue
		}

		dbUsersWithPrivileges := []entity.DatabaseUser{}
		for _, dbUsername := range dbUsernames {
			dbUsersWithPrivileges = append(
				dbUsersWithPrivileges, entity.NewDatabaseUser(
					dbUsername, dbName, dbType, []valueObject.DatabasePrivilege{"ALL PRIVILEGES"},
				),
			)
		}

		databaseEntities = append(
			databaseEntities, entity.NewDatabase(dbName, dbType, dbSize, dbUsersWithPrivileges),
		)
	}

	return databaseEntities, nil
}
