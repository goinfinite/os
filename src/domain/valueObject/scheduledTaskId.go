package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type ScheduledTaskId uint64

func NewScheduledTaskId(value interface{}) (
	scheduledTaskId ScheduledTaskId, err error,
) {
	uintValue, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return scheduledTaskId, errors.New("ScheduledTaskIdMustBeUint64")
	}

	return ScheduledTaskId(uintValue), nil
}

func (vo ScheduledTaskId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
