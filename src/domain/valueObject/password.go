package valueObject

import (
	"errors"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type Password string

func NewPassword(value interface{}) (password Password, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return password, errors.New("PasswordValueMustBeString")
	}

	if len(stringValue) < 6 {
		return password, errors.New("PasswordTooShort")
	}

	if len(stringValue) > 64 {
		return password, errors.New("PasswordTooLong")
	}

	return Password(stringValue), nil
}

func (vo Password) String() string {
	return string(vo)
}
