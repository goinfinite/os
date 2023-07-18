package valueObject

import (
	"errors"

	"golang.org/x/exp/slices"
)

type DatabasePrivilege string

var ValidDatabasePrivileges = []string{
	"openlitespeed",
	"nginx",
	"node",
	"mysql",
	"redis",
}

func NewDatabasePrivilege(value string) (DatabasePrivilege, error) {
	dp := DatabasePrivilege(value)
	if !dp.isValid() {
		return "", errors.New("InvalidDatabasePrivilege")
	}
	return dp, nil
}

func NewDatabasePrivilegePanic(value string) DatabasePrivilege {
	dp := DatabasePrivilege(value)
	if !dp.isValid() {
		panic("InvalidDatabasePrivilege")
	}
	return dp
}

func (dp DatabasePrivilege) isValid() bool {
	return slices.Contains(ValidDatabasePrivileges, dp.String())
}

func (dp DatabasePrivilege) String() string {
	return string(dp)
}
