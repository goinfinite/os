package valueObject

import (
	"errors"
	"strconv"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type CronId uint64

func NewCronId(value interface{}) (cronId CronId, err error) {
	uint64Value, err := tkVoUtil.InterfaceToUint64(value)
	if err != nil {
		return cronId, errors.New("CronIdMustBeUint64")
	}

	return CronId(uint64Value), nil
}

func (vo CronId) Uint64() uint64 {
	return uint64(vo)
}

func (vo CronId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
