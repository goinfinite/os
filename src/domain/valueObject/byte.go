package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type Byte int64

func NewByte(value interface{}) (byteValue Byte, err error) {
	uintValue, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return byteValue, errors.New("ByteMustBeUint64")
	}

	return Byte(uintValue), nil
}

func (vo Byte) Int64() int64 {
	return int64(vo)
}

func (vo Byte) ToKiB() int64 {
	return vo.Int64() / 1024
}

func (vo Byte) ToMiB() int64 {
	return vo.ToKiB() / 1024
}

func (vo Byte) ToGiB() int64 {
	return vo.ToMiB() / 1024
}

func (vo Byte) ToTiB() int64 {
	return vo.ToGiB() / 1024
}

func (vo Byte) String() string {
	return strconv.FormatInt(int64(vo), 10)
}
