package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type GroupId uint64

func NewGroupId(value interface{}) (groupId GroupId, err error) {
	uintValue, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return groupId, errors.New("GroupIdMustBeUint64")
	}

	return GroupId(uintValue), nil
}

func (vo GroupId) Uint64() uint64 {
	return uint64(vo)
}

func (vo GroupId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
