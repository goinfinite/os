package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type AccountId uint

func NewAccountId(value interface{}) (accountId AccountId, err error) {
	uintValue, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return accountId, errors.New("AccountIdMustBeUint")
	}

	return AccountId(uintValue), nil
}

func (vo AccountId) Uint() uint {
	return uint(vo)
}

func (vo AccountId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
