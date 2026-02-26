package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

const databaseNameRegex string = `^\w[\w-]{1,30}\w$`

type DatabaseName string

func NewDatabaseName(value interface{}) (dbName DatabaseName, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
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
