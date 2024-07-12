package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const scheduledTaskNameRegex string = `^[a-zA-Z][\w\-]{1,256}[\w\-\ ]{0,512}$`

type ScheduledTaskName string

func NewScheduledTaskName(value interface{}) (ScheduledTaskName, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ScheduledTaskNameMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)

	re := regexp.MustCompile(scheduledTaskNameRegex)
	isValid := re.MatchString(stringValue)
	if !isValid {
		return "", errors.New("InvalidScheduledTaskName")
	}

	return ScheduledTaskName(stringValue), nil
}

func (vo ScheduledTaskName) String() string {
	return string(vo)
}
