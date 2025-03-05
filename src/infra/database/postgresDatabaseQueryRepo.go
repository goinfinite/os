package databaseInfra

import (
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
)

type PostgresDatabaseQueryRepo struct {
}

func PostgresqlCmd(cmd string, dbName *string) (string, error) {
	psqlArgs := []string{"-U", "postgres", "-tAc", cmd}

	if dbName != nil {
		psqlDbToConnect := []string{"-d", *dbName}
		psqlArgs = append(psqlArgs, psqlDbToConnect...)
	}

	return infraHelper.RunCmd(infraHelper.RunCmdConfigs{
		Command: "psql",
		Args:    psqlArgs,
	})
}

func (repo PostgresDatabaseQueryRepo) getDatabaseNames() ([]valueObject.DatabaseName, error) {
	var dbNameList []valueObject.DatabaseName

	rawDbNameList, err := PostgresqlCmd("SELECT datname FROM pg_database", nil)
	if err != nil {
		return dbNameList, errors.New("GetDatabaseNamesError: " + err.Error())
	}

	rawDbNameListSlice := strings.Split(rawDbNameList, "\n")
	dbExcludeRegex := "^(postgres|template1|template0)$"
	compiledDbExcludeRegex := regexp.MustCompile(dbExcludeRegex)
	for _, rawDbName := range rawDbNameListSlice {
		if compiledDbExcludeRegex.MatchString(rawDbName) {
			continue
		}

		dbName, err := valueObject.NewDatabaseName(rawDbName)
		if err != nil {
			log.Printf("%s: %s", err.Error(), rawDbName)
			continue
		}

		dbNameList = append(dbNameList, dbName)
	}

	return dbNameList, nil
}

func (repo PostgresDatabaseQueryRepo) getDatabaseSize(
	dbName valueObject.DatabaseName,
) (valueObject.Byte, error) {
	rawDbSize, err := PostgresqlCmd(
		"SELECT pg_database_size('"+dbName.String()+"')",
		nil,
	)
	if err != nil {
		return 0, errors.New("GetDatabaseSizeError: " + err.Error())
	}

	return valueObject.NewByte(rawDbSize)
}

func (repo PostgresDatabaseQueryRepo) getDatabaseUsernames(
	dbName valueObject.DatabaseName,
) ([]valueObject.DatabaseUsername, error) {
	dbUsernameList := []valueObject.DatabaseUsername{}

	rawDbUsersPrivs, err := PostgresqlCmd(
		"SELECT datacl FROM pg_database WHERE datname = '"+dbName.String()+"'",
		nil,
	)
	if err != nil {
		return dbUsernameList, errors.New("GetDatabaseUserError: " + err.Error())
	}

	compiledDbUsersPrivsRegex := regexp.MustCompile(`(\w+)=`)
	rawDbUsersMatches := compiledDbUsersPrivsRegex.FindAllStringSubmatch(rawDbUsersPrivs, -1)

	if len(rawDbUsersMatches) == 0 {
		return dbUsernameList, nil
	}

	defaultDbUser := "postgres"
	for _, rawDbUserMatch := range rawDbUsersMatches {
		if len(rawDbUserMatch) < 2 {
			continue
		}

		rawDbUser := rawDbUserMatch[1]
		if rawDbUser == defaultDbUser {
			continue
		}

		dbUser, err := valueObject.NewDatabaseUsername(rawDbUser)
		if err != nil {
			log.Printf("%s: %s", err.Error(), rawDbUser)
			continue
		}

		dbUsernameList = append(dbUsernameList, dbUser)
	}

	return dbUsernameList, nil
}

func (repo PostgresDatabaseQueryRepo) Read() ([]entity.Database, error) {
	var databases []entity.Database

	dbNames, err := repo.getDatabaseNames()
	if err != nil {
		return databases, errors.New("GetDatabaseNamesError: " + err.Error())
	}
	dbType, _ := valueObject.NewDatabaseType("postgresql")

	for _, dbName := range dbNames {
		dbSize, err := repo.getDatabaseSize(dbName)
		if err != nil {
			dbSize, _ = valueObject.NewByte(0)
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

func (repo PostgresDatabaseQueryRepo) ReadDatabaseNamesByUser(
	dbUser valueObject.DatabaseUsername,
) ([]valueObject.DatabaseName, error) {
	dbNames := []valueObject.DatabaseName{}

	rawDbNames, err := PostgresqlCmd(
		"SELECT datname FROM pg_database WHERE array_to_string(datacl, '') LIKE '%"+
			dbUser.String()+"%'",
		nil,
	)
	if err != nil {
		return dbNames, errors.New("GetUserDatabaseNamesError: " + err.Error())
	}

	rawDbNamesSlice := strings.Split(rawDbNames, "\n")
	if len(rawDbNamesSlice) == 0 {
		return dbNames, nil
	}

	for _, rawDbName := range rawDbNamesSlice {
		if len(rawDbName) == 0 {
			continue
		}

		dbName, err := valueObject.NewDatabaseName(rawDbName)
		if err != nil {
			log.Printf("%s: %s", err.Error(), rawDbName)
			continue
		}

		dbNames = append(dbNames, dbName)
	}

	return dbNames, nil
}
