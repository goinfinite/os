package valueObject

import (
	"errors"
	"strconv"
)

type UserId int64

func NewUserIdFromString(value string) (UserId, error) {
	accId, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errors.New("InvalidUserId")
	}
	return UserId(accId), nil
}

func NewUserIdFromStringPanic(value string) UserId {
	accId, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic("InvalidUserId")
	}
	return UserId(accId)
}

func (uid UserId) GetUserId() int64 {
	return int64(uid)
}

func (uid UserId) String() string {
	return strconv.FormatInt(int64(uid), 10)
}
