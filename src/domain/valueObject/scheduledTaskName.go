package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const scheduledTaskNameRegex string = `^[a-zA-Z][\w\-]{1,256}[\w\-\ ]{0,512}$`

type ScheduledTaskName string

func NewScheduledTaskName(value interface{}) (ScheduledTaskName, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ScheduledTaskNameMustBeString")
	}

	re := regexp.MustCompile(scheduledTaskNameRegex)
	if !re.MatchString(stringValue) {
		return "", errors.New("InvalidScheduledTaskName")
	}

	return ScheduledTaskName(stringValue), nil
}

func (vo ScheduledTaskName) String() string {
	return string(vo)
}
