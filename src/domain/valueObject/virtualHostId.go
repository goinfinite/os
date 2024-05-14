package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type VirtualHostId uint

func NewVirtualHostId(value interface{}) (VirtualHostId, error) {
	valueUint, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidVirtualHostId")
	}

	vo := VirtualHostId(valueUint)
	if !vo.isValid() {
		return 0, errors.New("InvalidVirtualHostId")
	}

	return vo, nil
}

func NewVirtualHostIdPanic(value interface{}) VirtualHostId {
	vo, err := NewVirtualHostId(value)
	if err != nil {
		panic(err)
	}

	return vo
}

func (vo VirtualHostId) isValid() bool {
	return vo >= 1 && vo <= 999999999999
}

func (vo VirtualHostId) Get() uint {
	return uint(vo)
}

func (vo VirtualHostId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
