package valueObject

import (
	"errors"
	"io/fs"
	"regexp"
	"strconv"
)

var unixFilePermissionsRegex = regexp.MustCompile(`^[0-7]{3,4}$`)

type UnixFilePermissions string

/**
 * The "interfaceToUint" helper was not used due to the problem of octal
 * base vs decimal base in file permissions in C-like language.
 */
func NewUnixFilePermissions(value any) (
	unixFilePermission UnixFilePermissions, err error,
) {
	stringValue, assertOk := value.(string)
	if !assertOk {
		return unixFilePermission, errors.New("UnixFilePermissionsMustBeString")
	}

	if !unixFilePermissionsRegex.MatchString(stringValue) {
		return unixFilePermission, errors.New("InvalidUnixFilePermissions")
	}

	return UnixFilePermissions(stringValue), nil
}

func NewUnixFileDefaultPermissions() UnixFilePermissions {
	return UnixFilePermissions("644")
}

func NewUnixDirDefaultPermissions() UnixFilePermissions {
	return UnixFilePermissions("755")
}

func (vo UnixFilePermissions) GetFileMode() fs.FileMode {
	intValue, _ := strconv.ParseInt(string(vo), 8, 64)
	return fs.FileMode(intValue)
}

func (vo UnixFilePermissions) String() string {
	return string(vo)
}
