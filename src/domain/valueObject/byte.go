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

func (b Byte) Int64() int64 {
	return int64(b)
}

func (b Byte) ToKiB() int64 {
	return b.Int64() / 1024
}

func (b Byte) ToMiB() int64 {
	return b.ToKiB() / 1024
}

func (b Byte) ToGiB() int64 {
	return b.ToMiB() / 1024
}

func (b Byte) ToTiB() int64 {
	return b.ToGiB() / 1024
}

func (vo Byte) String() string {
	return strconv.FormatInt(int64(vo), 10)
}
