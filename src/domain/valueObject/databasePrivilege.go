package valueObject

import (
	"errors"
	"slices"
	"strings"
)

type DatabasePrivilege string

var ValidDatabasePrivileges = []string{
	"ALL PRIVILEGES",
	"ALL",
	"ALTER ROUTINE",
	"ALTER SYSTEM",
	"ALTER",
	"BYPASSRLS",
	"CONNECT",
	"CREATE DOMAIN",
	"CREATE FUNCTION",
	"CREATE GROUP",
	"CREATE INDEX",
	"CREATE LANGUAGE",
	"CREATE PROCEDURE",
	"CREATE ROLE",
	"CREATE ROUTINE",
	"CREATE SCHEMA",
	"CREATE TABLE",
	"CREATE TEMP",
	"CREATE TEMPORARY TABLES",
	"CREATE TRIGGER",
	"CREATE TYPE",
	"CREATE USER",
	"CREATE VIEW",
	"CREATE",
	"CREATEDB",
	"CREATEROLE",
	"DELETE HISTORY",
	"DELETE",
	"DROP",
	"EVENT",
	"EXECUTE",
	"FILE",
	"INDEX",
	"INSERT",
	"LOCK TABLES",
	"PASSWORDADMIN",
	"PROCESS",
	"PROXY",
	"REFERENCES",
	"RELOAD",
	"REPLICATION CLIENT",
	"REPLICATION SLAVE",
	"REPLICATION",
	"SELECT",
	"SET",
	"SHOW VIEW",
	"SHUTDOWN",
	"SUPER",
	"SUPERUSER",
	"TEMPORARY",
	"TRIGGER",
	"TRUNCATE",
	"UPDATE",
	"USAGE",
}

func NewDatabasePrivilege(value string) (DatabasePrivilege, error) {
	value = strings.ReplaceAll(value, "-", " ")

	dp := DatabasePrivilege(strings.ToUpper(value))
	if !dp.isValid() {
		return "", errors.New("InvalidDatabasePrivilege")
	}
	return dp, nil
}

func NewDatabasePrivilegePanic(value string) DatabasePrivilege {
	dp, err := NewDatabasePrivilege(value)
	if err != nil {
		panic(err)
	}
	return dp
}

func (dp DatabasePrivilege) isValid() bool {
	return slices.Contains(ValidDatabasePrivileges, dp.String())
}

func (dp DatabasePrivilege) String() string {
	return string(dp)
}
