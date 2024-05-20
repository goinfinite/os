package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type UnixUid int64

func NewUnixUid(value interface{}) (UnixUid, error) {
	unixUidInt, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return 0, errors.New("InvalidUnixUid")
	}

	unixUid := UnixUid(unixUidInt)
	if !unixUid.isValid() {
		return 0, errors.New("InvalidUnixUid")
	}

	return unixUid, nil
}

func NewUnixUidPanic(value int) UnixUid {
	unixUid, err := NewUnixUid(value)
	if err != nil {
		panic(err)
	}
	return unixUid
}

func (unixUid UnixUid) isValid() bool {
	return unixUid >= 0 && unixUid <= 2147483647
}

func (gid UnixUid) Get() int64 {
	return int64(gid)
}

func (gid UnixUid) String() string {
	return strconv.FormatInt(int64(gid), 10)
}
