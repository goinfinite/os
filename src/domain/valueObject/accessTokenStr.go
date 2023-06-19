package valueObject

import (
	"errors"
	"regexp"
)

type AccessTokenStr string

func NewAccessTokenStr(value string) (AccessTokenStr, error) {
	ats := AccessTokenStr(value)
	if !ats.isValid() {
		return "", errors.New("InvalidAccessTokenStr")
	}
	return ats, nil
}

func NewAccessTokenStrPanic(value string) AccessTokenStr {
	ats := AccessTokenStr(value)
	if !ats.isValid() {
		panic("InvalidAccessTokenStr")
	}
	return ats
}

func (ats AccessTokenStr) isValid() bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9\-_=+/.]{22,444}$`)
	return re.MatchString(string(ats))
}

func (ats AccessTokenStr) String() string {
	return string(ats)
}
