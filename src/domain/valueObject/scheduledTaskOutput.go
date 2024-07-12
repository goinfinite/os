package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type ScheduledTaskOutput string

func NewScheduledTaskOutput(value interface{}) (ScheduledTaskOutput, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ScheduledTaskOutputMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)
	valueLength := len(stringValue)
	if valueLength > 2048 {
		stringValue = stringValue[:2048]
	}

	return ScheduledTaskOutput(stringValue), nil
}

func (vo ScheduledTaskOutput) String() string {
	return string(vo)
}
