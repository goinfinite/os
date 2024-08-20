package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type NetworkPort uint16

func NewNetworkPort(value interface{}) (networkPort NetworkPort, err error) {
	uintValue, err := voHelper.InterfaceToUint16(value)
	if err != nil {
		return networkPort, errors.New("NetworkPortMustBeUint16")
	}

	return NetworkPort(uintValue), nil
}

func (vo NetworkPort) Uint() uint {
	return uint(vo)
}

func (vo NetworkPort) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
