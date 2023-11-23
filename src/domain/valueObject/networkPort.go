package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type NetworkPort uint64

func NewNetworkPort(value interface{}) (NetworkPort, error) {
	np, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidNetworkPort")
	}

	return NetworkPort(np), nil
}

func NewNetworkPortPanic(value interface{}) NetworkPort {
	np, err := NewNetworkPort(value)
	if err != nil {
		panic(err)
	}
	return np
}

func (np NetworkPort) Get() uint64 {
	return uint64(np)
}

func (np NetworkPort) String() string {
	return strconv.FormatUint(uint64(np), 10)
}
