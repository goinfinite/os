package valueObject

import (
	"errors"
	"slices"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type ScheduledTaskStatus string

var ValidScheduledTaskStatuses = []string{
	"pending", "running", "completed", "failed", "cancelled", "timeout",
}

func NewScheduledTaskStatus(value interface{}) (
	scheduledTaskStatus ScheduledTaskStatus, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return scheduledTaskStatus, errors.New("ScheduledTaskStatusMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidScheduledTaskStatuses, stringValue) {
		return scheduledTaskStatus, errors.New("InvalidScheduledTaskStatus")
	}

	return ScheduledTaskStatus(stringValue), nil
}

func (vo ScheduledTaskStatus) String() string {
	return string(vo)
}
