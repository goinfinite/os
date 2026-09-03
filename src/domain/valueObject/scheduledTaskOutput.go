package valueObject

import (
	"errors"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

const maxScheduledTaskOutputBytes = 65536

type ScheduledTaskOutput string

func NewScheduledTaskOutput(value any) (
	scheduledTaskOutput ScheduledTaskOutput, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return scheduledTaskOutput, errors.New("ScheduledTaskOutputMustBeString")
	}

	if len(stringValue) > maxScheduledTaskOutputBytes {
		stringValue = stringValue[:maxScheduledTaskOutputBytes]
	}

	return ScheduledTaskOutput(stringValue), nil
}

func (vo ScheduledTaskOutput) String() string {
	return string(vo)
}
