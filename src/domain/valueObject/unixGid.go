package valueObject

import (
	"errors"
	"strconv"
)

type UnixGid int64

func NewUnixGid(value int) (UnixGid, error) {
	unixGid := UnixGid(value)
	if !unixGid.isValid() {
		return 0, errors.New("InvalidUnixGid")
	}

	return unixGid, nil
}

func NewUnixGidPanic(value int) UnixGid {
	unixGid, err := NewUnixGid(value)
	if err != nil {
		panic("InvalidUnixGid")
	}
	return UnixGid(unixGid)
}

func (unixGid UnixGid) isValid() bool {
	return unixGid >= 0 && unixGid <= 2147483647
}

func (gid UnixGid) Get() int64 {
	return int64(gid)
}

func (gid UnixGid) String() string {
	return strconv.FormatInt(int64(gid), 10)
}
