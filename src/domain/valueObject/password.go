package valueObject

import (
	"errors"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type Password string

func NewPassword(value interface{}) (password Password, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("PasswordValueMustBeString")
	}

	valueLength := len(stringValue)
	if valueLength < 6 {
		return password, errors.New("PasswordTooShort")
	}

	if valueLength > 64 {
		return password, errors.New("PasswordTooLong")
	}

	return Password(stringValue), nil
}

// TODO: remove this constructor when no longer used.
func NewPasswordPanic(value interface{}) Password {
	pass, err := NewPassword(value)
	if err != nil {
		panic(err)
	}
	return pass
}

func (vo Password) String() string {
	return string(vo)
}
