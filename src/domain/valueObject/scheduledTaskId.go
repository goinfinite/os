package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type ScheduledTaskId uint

func NewScheduledTaskId(value interface{}) (ScheduledTaskId, error) {
	id, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidScheduledTaskId")
	}

	return ScheduledTaskId(id), nil
}

func (vo ScheduledTaskId) Read() uint64 {
	return uint64(vo)
}

func (vo ScheduledTaskId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
