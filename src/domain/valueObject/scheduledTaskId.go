package valueObject

import (
	"errors"
	"strconv"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type ScheduledTaskId uint64

func NewScheduledTaskId(value interface{}) (
	scheduledTaskId ScheduledTaskId, err error,
) {
	uintValue, err := tkVoUtil.InterfaceToUint64(value)
	if err != nil {
		return scheduledTaskId, errors.New("ScheduledTaskIdMustBeUint64")
	}

	return ScheduledTaskId(uintValue), nil
}

func (vo ScheduledTaskId) Uint64() uint64 {
	return uint64(vo)
}

func (vo ScheduledTaskId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
