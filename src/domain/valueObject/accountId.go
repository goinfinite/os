package valueObject

import (
	"errors"
	"strconv"
)

type AccountId int64

func NewAccountIdFromString(value string) (AccountId, error) {
	accId, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errors.New("InvalidAccountId")
	}
	return AccountId(accId), nil
}

func NewAccountIdFromStringPanic(value string) AccountId {
	accId, err := NewAccountIdFromString(value)
	if err != nil {
		panic(err)
	}
	return accId
}

func NewAccountIdFromFloat(value float64) (AccountId, error) {
	accId, err := strconv.ParseInt(
		strconv.FormatFloat(value, 'f', -1, 64), 10, 64,
	)
	if err != nil {
		return 0, errors.New("InvalidAccountId")
	}
	return AccountId(accId), nil
}

func NewAccountIdFromFloatPanic(value float64) AccountId {
	accId, err := NewAccountIdFromFloat(value)
	if err != nil {
		panic(err)
	}
	return AccountId(accId)
}

func (uid AccountId) Get() int64 {
	return int64(uid)
}

func (uid AccountId) String() string {
	return strconv.FormatInt(int64(uid), 10)
}
