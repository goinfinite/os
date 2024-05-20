package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type AccountId uint64

func NewAccountId(value interface{}) (AccountId, error) {
	accId, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return 0, errors.New("InvalidAccountId")
	}

	return AccountId(accId), nil
}

func NewAccountIdPanic(value interface{}) AccountId {
	accId, err := NewAccountId(value)
	if err != nil {
		panic(err)
	}
	return accId
}

func (id AccountId) Get() uint64 {
	return uint64(id)
}

func (id AccountId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}
