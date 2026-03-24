package valueObject

import (
	"errors"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type ScheduledTaskOutput string

func NewScheduledTaskOutput(value interface{}) (
	scheduledTaskOutput ScheduledTaskOutput, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return scheduledTaskOutput, errors.New("ScheduledTaskOutputMustBeString")
	}

	if len(stringValue) > 2048 {
		stringValue = stringValue[:2048]
	}

	return ScheduledTaskOutput(stringValue), nil
}

func (vo ScheduledTaskOutput) String() string {
	return string(vo)
}
