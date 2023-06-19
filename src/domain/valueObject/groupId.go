package valueObject

import (
	"errors"
	"strconv"
)

type GroupId int64

func NewGroupIdFromString(value string) (GroupId, error) {
	accId, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errors.New("InvalidGroupId")
	}
	return GroupId(accId), nil
}

func NewGroupIdFromStringPanic(value string) GroupId {
	accId, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic("InvalidGroupId")
	}
	return GroupId(accId)
}

func (gid GroupId) GetGroupId() int64 {
	return int64(gid)
}

func (gid GroupId) String() string {
	return strconv.FormatInt(int64(gid), 10)
}
