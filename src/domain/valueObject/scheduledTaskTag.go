package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const scheduledTaskTagRegex string = `^[a-zA-Z][\w\-]{1,256}$`

type ScheduledTaskTag string

func NewScheduledTaskTag(value interface{}) (ScheduledTaskTag, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("ScheduledTaskTagMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)
	stringValue = strings.ToLower(stringValue)

	re := regexp.MustCompile(scheduledTaskTagRegex)
	isValid := re.MatchString(stringValue)
	if !isValid {
		return "", errors.New("InvalidScheduledTaskTag")
	}

	return ScheduledTaskTag(stringValue), nil
}

func (vo ScheduledTaskTag) String() string {
	return string(vo)
}
