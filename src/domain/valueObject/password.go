package valueObject

import "errors"

type Password string

func NewPassword(value string) (Password, error) {
	pass := Password(value)
	if !pass.isValid() {
		return "", errors.New("InvalidPassword")
	}
	return pass, nil
}

func NewPasswordPanic(value string) Password {
	pass := Password(value)
	if !pass.isValid() {
		panic("InvalidPassword")
	}
	return pass
}

func (pass Password) isValid() bool {
	isTooShort := len(string(pass)) < 6
	isTooLong := len(string(pass)) > 64
	return !isTooShort && !isTooLong
}

func (pass Password) String() string {
	return string(pass)
}
