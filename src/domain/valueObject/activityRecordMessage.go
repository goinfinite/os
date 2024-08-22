package valueObject

import (
	"errors"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type ActivityRecordMessage string

func NewActivityRecordMessage(value interface{}) (
	activityRecordMessage ActivityRecordMessage, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return activityRecordMessage, errors.New("ActivityRecordMessageMustBeString")
	}

	if len(stringValue) > 2048 {
		stringValue = stringValue[:2048]
	}

	return ActivityRecordMessage(stringValue), nil
}

func (vo ActivityRecordMessage) String() string {
	return string(vo)
}
