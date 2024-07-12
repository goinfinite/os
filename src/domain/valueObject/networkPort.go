package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type NetworkPort uint

func NewNetworkPort(value interface{}) (NetworkPort, error) {
	np, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidNetworkPort")
	}

	return NetworkPort(np), nil
}

func (vo NetworkPort) Get() uint {
	return uint(vo)
}

func (vo NetworkPort) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
