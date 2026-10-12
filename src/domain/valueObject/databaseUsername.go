package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var databaseUsernameRegex = regexp.MustCompile(`^\w[\w-]+\w$`)

type DatabaseUsername string

func NewDatabaseUsername(value any) (
	dbUsername DatabaseUsername, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return dbUsername, errors.New("DatabaseUsernameMustBeString")
	}

	if !databaseUsernameRegex.MatchString(stringValue) {
		return dbUsername, errors.New("InvalidDatabaseUsername")
	}

	return DatabaseUsername(stringValue), nil
}

func (vo DatabaseUsername) String() string {
	return string(vo)
}
