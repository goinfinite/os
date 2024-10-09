package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const databaseUsernameRegex string = `^\w[\w-]+\w$`

type DatabaseUsername string

func NewDatabaseUsername(value interface{}) (
	dbUsername DatabaseUsername, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return dbUsername, errors.New("DatabaseUsernameMustBeString")
	}

	re := regexp.MustCompile(databaseUsernameRegex)
	if !re.MatchString(stringValue) {
		return dbUsername, errors.New("InvalidDatabaseUsername")
	}

	return DatabaseUsername(stringValue), nil
}

func (vo DatabaseUsername) String() string {
	return string(vo)
}
