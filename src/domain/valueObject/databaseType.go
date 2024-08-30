package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
	"golang.org/x/exp/maps"
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
	databaseTypes := maps.Keys(databaseTypesWithAliases)
	if slices.Contains(databaseTypes, value) {
		return value, nil
	}

	for _, databaseType := range databaseTypes {
		if !slices.Contains(databaseTypesWithAliases[databaseType], value) {
			continue
		}

		return databaseType, nil
	}

	return "", errors.New("InvalidDatabaseType")
}

func (vo DatabaseType) String() string {
	return string(vo)
}
