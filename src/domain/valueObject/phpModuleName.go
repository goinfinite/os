package valueObject

import (
	"errors"
	"slices"
	"strings"
)

var ValidPhpModuleNames = []string{
	"curl",
	"mysqli",
	"opcache",
	"apcu",
	"igbinary",
	"imagick",
	"imap",
	"intl",
	"ioncube",
	"ldap",
	"mailparse",
	"memcached",
	"mcrypt",
	"mongodb",
	"msgpack",
	"parallel",
	"pdo_mysql",
	"pdo_sqlite",
	"pear",
	"pgsql",
	"phalcon",
	"pspell",
	"redis",
	"snmp",
	"solr",
	"sqlite3",
	"sqlsrv",
	"ssh2",
	"swoole",
	"sybase",
	"tidy",
	"timezonedb",
	"yaml",
	"xdebug",
}

type PhpModuleName string

func NewPhpModuleName(value string) (PhpModuleName, error) {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)

	moduleName := PhpModuleName(value)
	if !moduleName.isValid() {
		return "", errors.New("InvalidPhpModuleName")
	}
	return moduleName, nil
}

func NewPhpModuleNamePanic(value string) PhpModuleName {
	moduleName, err := NewPhpModuleName(value)
	if err != nil {
		panic("InvalidPhpModuleName")
	}
	return moduleName
}

func (moduleName PhpModuleName) isValid() bool {
	return slices.Contains(ValidPhpModuleNames, moduleName.String())
}

func (moduleName PhpModuleName) String() string {
	return string(moduleName)
}
