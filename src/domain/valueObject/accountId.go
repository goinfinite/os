package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type AccountId uint64

func NewAccountId(value interface{}) (accountId AccountId, err error) {
	uintValue, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return accountId, errors.New("AccountIdMustBeUint64")
	}

	return AccountId(uintValue), nil
}

func (vo AccountId) Read() uint64 {
	return uint64(vo)
}

func (vo AccountId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
