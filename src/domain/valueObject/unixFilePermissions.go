package valueObject

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const unixFilePermissionsRegexExpression = `^[0-7]{1,4}$`

type UnixFilePermissions string

func NewUnixFilePermissions(value string) (UnixFilePermissions, error) {
	if len(value) < 1 {
		return "", errors.New("InvalidUnixFilePermissions")
	}

	paddedValue := fmt.Sprintf("%04s", value)

	unixFilePermissions := UnixFilePermissions(paddedValue)
	if !unixFilePermissions.isValid() {
		return "", errors.New("InvalidUnixFilePermissions")
	}

	return unixFilePermissions, nil
}

func NewUnixFilePermissionsPanic(value string) UnixFilePermissions {
	unixFilePermissions, err := NewUnixFilePermissions(value)
	if err != nil {
		panic(err)
	}
	return UnixFilePermissions(unixFilePermissions)
}

func (unixFilePermissions UnixFilePermissions) isValid() bool {
	unixFilePermissionsRegex := regexp.MustCompile(unixFilePermissionsRegexExpression)
	return unixFilePermissionsRegex.MatchString(string(unixFilePermissions))
}

func NewUnixFilePermissionsFromInt(value interface{}) (UnixFilePermissions, error) {
	intValue, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return "", errors.New("InvalidUnixFilePermissions")
	}

	valueStr := strconv.FormatUint(intValue, 10)

	return NewUnixFilePermissions(valueStr)
}

func NewUnixFilePermissionsFromIntPanic(value interface{}) UnixFilePermissions {
	unixFilePermissions, err := NewUnixFilePermissionsFromInt(value)
	if err != nil {
		panic(err)
	}
	return UnixFilePermissions(unixFilePermissions)
}

func (unixFilePermission UnixFilePermissions) String() string {
	return string(unixFilePermission)
}
