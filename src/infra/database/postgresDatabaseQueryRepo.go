package databaseInfra

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

type PostgresDatabaseQueryRepo struct {
}

func PostgresqlCmd(cmd string, dbName *string) (string, error) {
	psqlArgs := []string{"-U", "postgres", "-tAc", cmd}

	if dbName != nil {
		psqlDbToConnect := []string{"-d", *dbName}
		psqlArgs = append(psqlArgs, psqlDbToConnect...)
	}

	return infraHelper.RunCmd("psql", psqlArgs...)
}

func (repo PostgresDatabaseQueryRepo) getDatabaseNames() ([]valueObject.DatabaseName, error) {
	var dbNameList []valueObject.DatabaseName

	dbNameListStr, err := PostgresqlCmd("SELECT datname FROM pg_database", nil)
	if err != nil {
		return dbNameList, errors.New("GetDatabaseNamesError: " + err.Error())
	}

	dbNameListSlice := strings.Split(dbNameListStr, "\n")
	dbExcludeRegex := "^(postgres|template1|template0)$"
	compiledDbExcludeRegex := regexp.MustCompile(dbExcludeRegex)
	for _, dbNameStr := range dbNameListSlice {
		if compiledDbExcludeRegex.MatchString(dbNameStr) {
			continue
		}

		dbName, err := valueObject.NewDatabaseName(dbNameStr)
		if err != nil {
			log.Printf("%s: %s", err.Error(), dbNameStr)
			continue
		}

		dbNameList = append(dbNameList, dbName)
	}

	return dbNameList, nil
}

func (repo PostgresDatabaseQueryRepo) getDatabaseSize(
	dbName valueObject.DatabaseName,
) (valueObject.Byte, error) {
	dbSizeStr, err := PostgresqlCmd(
		"SELECT pg_database_size('"+dbName.String()+"')",
		nil,
	)
	if err != nil {
		return 0, errors.New("GetDatabaseSizeError: " + err.Error())
	}

	dbSizeInBytes, err := strconv.ParseInt(dbSizeStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return valueObject.Byte(dbSizeInBytes), nil
}

func (repo PostgresDatabaseQueryRepo) getDatabaseUsernames(
	dbName valueObject.DatabaseName,
) ([]valueObject.DatabaseUsername, error) {
	dbUsernameList := []valueObject.DatabaseUsername{}

	dbUsersPrivs, err := PostgresqlCmd(
		"SELECT datacl FROM pg_database WHERE datname = '"+dbName.String()+"'",
		nil,
	)
	if err != nil {
		return dbUsernameList, errors.New("GetDatabaseUserError: " + err.Error())
	}

	compiledDbUsersPrivsRegex := regexp.MustCompile(`(\w+)=`)
	dbUsersMatches := compiledDbUsersPrivsRegex.FindAllStringSubmatch(dbUsersPrivs, -1)

	if len(dbUsersMatches) == 0 {
		return dbUsernameList, nil
	}

	defaultDbUser := "postgres"
	for _, dbUserMatch := range dbUsersMatches {
		if len(dbUserMatch) < 2 {
			continue
		}

		dbUserStr := dbUserMatch[1]
		if dbUserStr == defaultDbUser {
			continue
		}

		dbUser, err := valueObject.NewDatabaseUsername(dbUserStr)
		if err != nil {
			log.Printf("%s: %s", err.Error(), dbUserStr)
			continue
		}

		dbUsernameList = append(dbUsernameList, dbUser)
	}

	return dbUsernameList, nil
}

func (repo PostgresDatabaseQueryRepo) Get() ([]entity.Database, error) {
	var databases []entity.Database

	dbNames, err := repo.getDatabaseNames()
	if err != nil {
		return databases, errors.New("GetDatabaseNamesError: " + err.Error())
	}
	dbType, _ := valueObject.NewDatabaseType("postgresql")

	for _, dbName := range dbNames {
		dbSize, err := repo.getDatabaseSize(dbName)
		if err != nil {
			dbSize = valueObject.Byte(0)
		}

		dbUsernames, err := repo.getDatabaseUsernames(dbName)
		if err != nil {
			log.Printf("GetDatabaseUsersError (%s): %s", dbName.String(), err.Error())
		}

		dbUsersWithPrivileges := []entity.DatabaseUser{}
		for _, dbUsername := range dbUsernames {
			dbUsersWithPrivileges = append(
				dbUsersWithPrivileges,
				entity.NewDatabaseUser(
					dbUsername,
					dbName,
					dbType,
					[]valueObject.DatabasePrivilege{"ALL PRIVILEGES"},
				),
			)
		}

		databases = append(
			databases,
			entity.NewDatabase(
				dbName,
				dbType,
				dbSize,
				dbUsersWithPrivileges,
			),
		)
	}

	return databases, nil
}

func (repo PostgresDatabaseQueryRepo) UserExists(
	dbUser valueObject.DatabaseUsername,
) bool {
	userExists, err := PostgresqlCmd(
		"SELECT 1 FROM pg_user WHERE usename='"+dbUser.String()+"'",
		nil,
	)
	if err != nil {
		return false
	}

	return userExists == "1"
}

func (repo PostgresDatabaseQueryRepo) GetDatabaseNamesByUser(
	dbUser valueObject.DatabaseUsername,
) ([]valueObject.DatabaseName, error) {
	dbNames := []valueObject.DatabaseName{}

	dbNamesStr, err := PostgresqlCmd(
		"SELECT datname FROM pg_database WHERE array_to_string(datacl, '') LIKE '%"+
			dbUser.String()+"%'",
		nil,
	)
	if err != nil {
		return dbNames, errors.New("GetUserDatabaseNamesError: " + err.Error())
	}

	dbNamesStrSlice := strings.Split(dbNamesStr, "\n")
	if len(dbNamesStrSlice) == 0 {
		return dbNames, nil
	}

	for _, dbNameStr := range dbNamesStrSlice {
		if len(dbNameStr) == 0 {
			continue
		}

		dbName, err := valueObject.NewDatabaseName(dbNameStr)
		if err != nil {
			log.Printf("%s: %s", err.Error(), dbNameStr)
			continue
		}

		dbNames = append(dbNames, dbName)
	}

	return dbNames, nil
}
