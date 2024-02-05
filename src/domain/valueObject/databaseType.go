package valueObject

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type DatabaseType string

const databaseTypeRegExp string = `^[a-z0-9\.\_\-]{1,64}$`

var databaseTypesWithAliases = map[string][]string{
	"mysql":    {"mariadb", "percona"},
	"postgres": {"postgresql"},
}

func NewDatabaseType(value string) (DatabaseType, error) {
	value = strings.ToLower(value)
	value = databaseTypeAdapter(value)

	dt := DatabaseType(value)
	if !dt.isValid() {
		return "", errors.New("InvalidDatabaseType")
	}
	return dt, nil
}

func NewDatabaseTypePanic(value string) DatabaseType {
	dt, err := NewDatabaseType(value)
	if err != nil {
		panic(err.Error())
	}
	return dt
}

func (dt DatabaseType) isValid() bool {
	databaseTypeRegex := regexp.MustCompile(databaseTypeRegExp)
	return databaseTypeRegex.MatchString(string(dt))
}

func databaseTypeAdapter(value string) string {
	databaseTypes := maps.Keys(databaseTypesWithAliases)
	for _, databaseType := range databaseTypes {
		if !slices.Contains(databaseTypesWithAliases[databaseType], value) {
			continue
		}

		return databaseType
	}

	return value
}

func (dt DatabaseType) String() string {
	return string(dt)
}
