package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const hashRegex string = `^\w{6,256}$`

type Hash string

func NewHash(value interface{}) (hash Hash, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return hash, errors.New("HashMustBeString")
	}

	re := regexp.MustCompile(hashRegex)
	if !re.MatchString(stringValue) {
		return hash, errors.New("InvalidHash")
	}

	return Hash(stringValue), nil
}

func (vo Hash) String() string {
	return string(vo)
}
