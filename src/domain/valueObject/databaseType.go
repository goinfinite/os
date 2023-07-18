package valueObject

import "errors"

type DatabaseType string

const (
	mysql    DatabaseType = "mysql"
	postgres DatabaseType = "postgres"
)

func NewDatabaseType(value string) (DatabaseType, error) {
	dt := DatabaseType(value)
	if !dt.isValid() {
		return "", errors.New("InvalidDatabaseType")
	}
	return dt, nil
}

func NewDatabaseTypePanic(value string) DatabaseType {
	dt := DatabaseType(value)
	if !dt.isValid() {
		panic("InvalidDatabaseType")
	}
	return dt
}

func (dt DatabaseType) isValid() bool {
	switch dt {
	case mysql, postgres:
		return true
	default:
		return false
	}
}

func (dt DatabaseType) String() string {
	return string(dt)
}
