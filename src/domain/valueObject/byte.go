package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type Byte int64

func NewByte(value interface{}) (Byte, error) {
	intValue, err := voHelper.InterfaceToInt64(value)
	if err != nil {
		return 0, errors.New("InvalidByte")
	}

	return Byte(intValue), nil
}

func NewMebibyte(value interface{}) (Byte, error) {
	intValue, err := voHelper.InterfaceToInt64(value)
	if err != nil {
		return 0, errors.New("InvalidMebibytes")
	}

	return Byte(intValue * 1048576), nil
}

func NewGibibyte(value interface{}) (Byte, error) {
	intValue, err := voHelper.InterfaceToInt64(value)
	if err != nil {
		return 0, errors.New("InvalidGibibytes")
	}

	return Byte(intValue * 1073741824), nil
}

func (vo Byte) Int64() int64 {
	return int64(vo)
}

func (vo Byte) Uint64() uint64 {
	return uint64(vo)
}

func (vo Byte) ToKiB() int64 {
	return vo.Int64() / 1024
}

func (vo Byte) ToMiB() int64 {
	return vo.Int64() / 1048576
}

func (vo Byte) ToGiB() int64 {
	return vo.Int64() / 1073741824
}

func (vo Byte) ToTiB() int64 {
	return vo.Int64() / 1099511627776
}

func (vo Byte) String() string {
	return strconv.FormatInt(int64(vo), 10)
}

func (vo Byte) StringWithSuffix() string {
	voInt64 := vo.Int64()
	switch {
	case voInt64 < 1048576:
		return strconv.FormatInt(vo.ToKiB(), 10) + " KiB"
	case voInt64 < 1073741824:
		return strconv.FormatInt(vo.ToMiB(), 10) + " MiB"
	case voInt64 < 1099511627776:
		return strconv.FormatInt(vo.ToGiB(), 10) + " GiB"
	case voInt64 < 1125899906842624:
		return strconv.FormatInt(vo.ToTiB(), 10) + " TiB"
	default:
		return strconv.FormatInt(voInt64, 10) + " B"
	}
}
