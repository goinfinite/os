package valueObject

import (
	"errors"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type ScheduledTaskOutput string

func NewScheduledTaskOutput(value interface{}) (ScheduledTaskOutput, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ScheduledTaskOutputMustBeString")
	}

	valueLength := len(stringValue)
	if valueLength > 2048 {
		stringValue = stringValue[:2048]
	}

	return ScheduledTaskOutput(stringValue), nil
}

func (vo ScheduledTaskOutput) String() string {
	return string(vo)
}
