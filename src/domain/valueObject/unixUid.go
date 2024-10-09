package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type UnixUid int64

func NewUnixUid(value interface{}) (unixUid UnixUid, err error) {
	intValue, err := voHelper.InterfaceToInt64(value)
	if err != nil {
		return unixUid, errors.New("InvalidUnixUid")
	}

	return UnixUid(intValue), nil
}

func (vo UnixUid) Int64() int64 {
	return int64(vo)
}

func (vo UnixUid) String() string {
	return strconv.FormatInt(int64(vo), 10)
}
