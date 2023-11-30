package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type GroupId int64

func NewGroupId(value interface{}) (GroupId, error) {
	groupIdInt, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidGroupId")
	}

	groupId := GroupId(groupIdInt)
	if !groupId.isValid() {
		return 0, errors.New("InvalidGroupId")
	}

	return groupId, nil
}

func NewGroupIdPanic(value interface{}) GroupId {
	groupId, err := NewGroupId(value)
	if err != nil {
		panic(err)
	}

	return groupId
}

func (gid GroupId) isValid() bool {
	return gid >= 0 && gid <= 999999999999
}

func (gid GroupId) Get() int64 {
	return int64(gid)
}

func (gid GroupId) String() string {
	return strconv.FormatInt(int64(gid), 10)
}
