package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const databaseNameRegex string = `^\w[\w-]{1,30}\w$`

type DatabaseName string

func NewDatabaseName(value interface{}) (dbName DatabaseName, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return dbName, errors.New("DatabaseNameMustBeString")
	}

	re := regexp.MustCompile(databaseNameRegex)
	if !re.MatchString(stringValue) {
		return dbName, errors.New("InvalidDatabaseName")
	}

	return DatabaseName(stringValue), nil
}

func (vo DatabaseName) String() string {
	return string(vo)
}
