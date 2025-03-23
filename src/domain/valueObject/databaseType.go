package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type DatabaseType string

var databaseTypesWithAliases = map[string][]string{
	"mariadb":    {"mysql", "percona"},
	"postgresql": {"postgres"},
}

func NewDatabaseType(value interface{}) (dbType DatabaseType, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return dbType, errors.New("DatabaseTypeMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	stringValue, err = databaseTypeAdapter(stringValue)
	if err != nil {
		return dbType, errors.New("InvalidDatabaseType")
	}

	return DatabaseType(stringValue), nil
}

func databaseTypeAdapter(value string) (string, error) {
	for exactName, aliases := range databaseTypesWithAliases {
		if exactName == value {
			return exactName, nil
		}

		if slices.Contains(aliases, value) {
			return exactName, nil
		}
	}

	return "", errors.New("InvalidDatabaseType")
}

func (vo DatabaseType) String() string {
	return string(vo)
}
