package valueObject

import (
	"errors"
	"regexp"
)

const databaseNameRegex string = `^[0-9a-zA-Z_-]{3,32}$`

type DatabaseName string

func NewDatabaseName(value string) (DatabaseName, error) {
	dbName := DatabaseName(value)
	if !dbName.isValid() {
		return "", errors.New("InvalidDatabaseName")
	}
	return dbName, nil
}

func NewDatabaseNamePanic(value string) DatabaseName {
	dbName := DatabaseName(value)
	if !dbName.isValid() {
		panic("InvalidDatabaseName")
	}
	return dbName
}

func (dbName DatabaseName) isValid() bool {
	re := regexp.MustCompile(databaseNameRegex)
	return re.MatchString(string(dbName))
}

func (dbName DatabaseName) String() string {
	return string(dbName)
}
