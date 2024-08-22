package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type ActivityRecordId uint64

func NewActivityRecordId(value interface{}) (
	activityRecordId ActivityRecordId, err error,
) {
	uintValue, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return activityRecordId, errors.New("ActivityRecordIdMustBeUint64")
	}

	return ActivityRecordId(uintValue), nil
}

func (vo ActivityRecordId) Uint64() uint64 {
	return uint64(vo)
}

func (vo ActivityRecordId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
