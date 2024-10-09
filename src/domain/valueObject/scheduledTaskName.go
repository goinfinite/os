package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const scheduledTaskNameRegex string = `^[a-zA-Z][\w\-]{1,256}[\w\-\ ]{0,512}$`

type ScheduledTaskName string

func NewScheduledTaskName(value interface{}) (
	scheduledTaskName ScheduledTaskName, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return scheduledTaskName, errors.New("ScheduledTaskNameMustBeString")
	}

	re := regexp.MustCompile(scheduledTaskNameRegex)
	if !re.MatchString(stringValue) {
		return scheduledTaskName, errors.New("InvalidScheduledTaskName")
	}

	return ScheduledTaskName(stringValue), nil
}

func (vo ScheduledTaskName) String() string {
	return string(vo)
}
