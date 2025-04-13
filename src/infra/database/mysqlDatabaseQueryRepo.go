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

type MysqlDatabaseQueryRepo struct {
}

func MysqlCmd(cmd string) (string, error) {
	return infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "mysql",
		Args: []string{
			"--defaults-file=/root/.my.cnf", "--skip-column-names", "--silent",
			"--execute", cmd,
		},
	})
}

func (repo MysqlDatabaseQueryRepo) readDatabaseNames() ([]valueObject.DatabaseName, error) {
	databaseNames := []valueObject.DatabaseName{}

	rawDatabaseNames, err := MysqlCmd("SHOW DATABASES")
	if err != nil {
		return databaseNames, errors.New("ReadDatabaseNamesError: " + err.Error())
	}

	dbExcludeRegex := regexp.MustCompile(
		`^(information_schema|mysql|performance_schema|sys)$`,
	)
	for rawDatabaseName := range strings.SplitSeq(rawDatabaseNames, "\n") {
		if dbExcludeRegex.MatchString(rawDatabaseName) {
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

func (repo MysqlDatabaseQueryRepo) readDatabaseSize(dbName valueObject.DatabaseName) (
	valueObject.Byte,
	error,
) {
	rawDatabaseSize, err := MysqlCmd(
		"SELECT SUM(data_length + index_length) FROM information_schema.TABLES WHERE table_schema = '" + dbName.String() + "'",
	)
	if err != nil {
		return 0, errors.New("ReadDatabaseSizeError: " + err.Error())
	}

	return valueObject.NewByte(rawDatabaseSize)
}

func (repo MysqlDatabaseQueryRepo) readDatabaseUsernames(
	dbName valueObject.DatabaseName,
) ([]valueObject.DatabaseUsername, error) {
	dbUsernames := []valueObject.DatabaseUsername{}

	rawDatabaseUsers, err := MysqlCmd(
		"SELECT User FROM mysql.db WHERE Db = '" + dbName.String() + "'",
	)
	if err != nil {
		return dbUsernames, errors.New("ReadDatabaseUserError: " + err.Error())
	}
	rawDatabaseUsers = strings.TrimSpace(rawDatabaseUsers)

	for rawDatabaseUsername := range strings.SplitSeq(rawDatabaseUsers, "\n") {
		if rawDatabaseUsername == "" {
			continue
		}

		dbUsername, err := valueObject.NewDatabaseUsername(rawDatabaseUsername)
		if err != nil {
			slog.Debug(err.Error(), slog.String("rawDbUser", rawDatabaseUsername))
			continue
		}
		dbUsernames = append(dbUsernames, dbUsername)
	}

	return dbUsernames, nil
}

func (repo MysqlDatabaseQueryRepo) readDatabaseUserPrivileges(
	dbName valueObject.DatabaseName,
	dbUser valueObject.DatabaseUsername,
) ([]valueObject.DatabasePrivilege, error) {
	dbUserPrivileges := []valueObject.DatabasePrivilege{}

	userGrantsStr, err := MysqlCmd(
		"SHOW GRANTS FOR '" + dbUser.String() + "'",
	)
	if err != nil {
		return dbUserPrivileges, errors.New(
			"ReadDatabaseUserPrivilegesError: " + err.Error(),
		)
	}

	grantsRegexp := regexp.MustCompile(
		`GRANT (?P<privs>.*) ON (?:\x60|'|")?` + dbName.String() + `(?:\x60|'|")?\.`,
	)
	for rawGrants := range strings.SplitSeq(userGrantsStr, "\n") {
		rawGrants = strings.TrimSpace(rawGrants)
		if !grantsRegexp.MatchString(rawGrants) {
			continue
		}

		rawGrantsRegexParts := grantsRegexp.FindStringSubmatch(rawGrants)
		if len(rawGrantsRegexParts) < 2 {
			continue
		}

		rawPrivileges := rawGrantsRegexParts[1]
		for rawSinglePrivilege := range strings.SplitSeq(rawPrivileges, ",") {
			privilege, err := valueObject.NewDatabasePrivilege(rawSinglePrivilege)
			if err != nil {
				slog.Debug(
					err.Error(),
					slog.String("rawPrivilege", rawSinglePrivilege),
					slog.String("dbName", dbName.String()),
					slog.String("dbUser", dbUser.String()),
				)
				continue
			}

			dbUserPrivileges = append(dbUserPrivileges, privilege)
		}
	}

	return dbUserPrivileges, nil
}

func (repo MysqlDatabaseQueryRepo) readAllDatabases() ([]entity.Database, error) {
	databaseEntities := []entity.Database{}

	dbNames, err := repo.readDatabaseNames()
	if err != nil {
		return databaseEntities, err
	}
	dbType, _ := valueObject.NewDatabaseType("mariadb")

	for _, dbName := range dbNames {
		dbSize, err := repo.readDatabaseSize(dbName)
		if err != nil {
			dbSize, _ = valueObject.NewByte(0)
		}

		dbUsernames, err := repo.readDatabaseUsernames(dbName)
		if err != nil {
			dbUsernames = []valueObject.DatabaseUsername{}
		}

		dbUsersWithPrivileges := []entity.DatabaseUser{}
		for _, dbUsername := range dbUsernames {
			dbUserPrivileges, err := repo.readDatabaseUserPrivileges(dbName, dbUsername)
			if err != nil {
				dbUserPrivileges = []valueObject.DatabasePrivilege{}
			}

			dbUsersWithPrivileges = append(
				dbUsersWithPrivileges,
				entity.NewDatabaseUser(dbUsername, dbName, dbType, dbUserPrivileges),
			)
		}

		databaseEntities = append(
			databaseEntities, entity.NewDatabase(dbName, dbType, dbSize, dbUsersWithPrivileges),
		)
	}

	return databaseEntities, nil
}
