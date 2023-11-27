package valueObject

import (
	"errors"
	"strconv"
)

type UnixUid int64

func NewUnixUid(value int) (UnixUid, error) {
	unixUid := UnixUid(value)
	if !unixUid.isValid() {
		return 0, errors.New("InvalidUnixUid")
	}

	return unixUid, nil
}

func NewUnixUidPanic(value int) UnixUid {
	unixUid, err := NewUnixUid(value)
	if err != nil {
		panic("InvalidUnixUid")
	}
	return UnixUid(unixUid)
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
