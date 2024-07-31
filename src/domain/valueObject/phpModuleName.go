package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
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

func NewPhpModuleName(value interface{}) (moduleName PhpModuleName, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return moduleName, errors.New("PhpModuleNameMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidPhpModuleNames, stringValue) {
		return moduleName, errors.New("InvalidPhpModuleName")
	}
	return PhpModuleName(stringValue), nil
}

func (vo PhpModuleName) String() string {
	return string(vo)
}
