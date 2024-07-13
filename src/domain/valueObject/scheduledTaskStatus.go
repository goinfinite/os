package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type ScheduledTaskStatus string

var ValidScheduledTaskStatuses = []string{
	"pending", "running", "completed", "failed", "cancelled", "timeout",
}

func NewScheduledTaskStatus(value interface{}) (ScheduledTaskStatus, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ScheduledTaskStatusMustBeString")
	}

	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidScheduledTaskStatuses, stringValue) {
		return "", errors.New("InvalidScheduledTaskStatus")
	}
	return ScheduledTaskStatus(stringValue), nil
}

func (vo ScheduledTaskStatus) String() string {
	return string(vo)
}
