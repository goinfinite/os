package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type AccountId uint64

func NewAccountId(value interface{}) (accountId AccountId, err error) {
	uintValue, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return accountId, errors.New("AccountIdMustBeUint")
	}

	return AccountId(uintValue), nil
}

func (vo AccountId) Uint64() uint64 {
	return uint64(vo)
}

func (vo AccountId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
