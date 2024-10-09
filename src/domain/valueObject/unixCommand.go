package valueObject

import (
	"errors"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type UnixCommand string

func NewUnixCommand(value interface{}) (unixCommand UnixCommand, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return unixCommand, errors.New("UnixCommandValueMustBeString")
	}

	if len(stringValue) < 2 {
		return unixCommand, errors.New("UnixCommandTooShort")
	}

	if len(stringValue) > 4096 {
		return unixCommand, errors.New("UnixCommandTooLong")
	}

	return UnixCommand(stringValue), nil
}

func (vo UnixCommand) String() string {
	return string(vo)
}
