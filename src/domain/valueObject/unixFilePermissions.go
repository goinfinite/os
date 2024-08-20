package valueObject

import (
	"errors"
	"io/fs"
	"regexp"
	"strconv"
)

const unixFilePermissionsRegexExpression = `^[0-7]{3,4}$`

type UnixFilePermissions string

/**
 * The "interfaceToUint" helper was not used due to the problem of octal
 * base vs decimal base in file permissions in C-like language.
 */
func NewUnixFilePermissions(value interface{}) (
	unixFilePermission UnixFilePermissions, err error,
) {
	stringValue, assertOk := value.(string)
	if !assertOk {
		return unixFilePermission, errors.New("UnixFilePermissionsMustBeString")
	}

	re := regexp.MustCompile(unixFilePermissionsRegexExpression)
	if !re.MatchString(stringValue) {
		return unixFilePermission, errors.New("InvalidUnixFilePermissions")
	}

	return UnixFilePermissions(stringValue), nil
}

func (vo UnixFilePermissions) GetFileMode() fs.FileMode {
	intValue, _ := strconv.ParseInt(string(vo), 8, 64)
	return fs.FileMode(intValue)
}

func (vo UnixFilePermissions) String() string {
	return string(vo)
}
