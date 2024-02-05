package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type DatabaseType string

var databaseTypesWithAliases = map[string][]string{
	"mysql":    {"mariadb", "percona"},
	"postgres": {"postgresql"},
}

func NewDatabaseType(value string) (DatabaseType, error) {
	value = strings.ToLower(value)
	value, err := databaseTypeAdapter(value)
	if err != nil {
		return "", errors.New("InvalidDatabaseType")
	}

	return DatabaseType(value), nil
}

func NewDatabaseTypePanic(value string) DatabaseType {
	dt, err := NewDatabaseType(value)
	if err != nil {
		panic(err.Error())
	}
	return dt
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

func (dt DatabaseType) String() string {
	return string(dt)
}
