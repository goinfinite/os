package valueObject

import (
	"errors"
	"regexp"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

const scheduledTaskTagRegex string = `^[a-zA-Z][\w\-]{1,256}$`

type ScheduledTaskTag string

func NewScheduledTaskTag(value interface{}) (
	scheduledTaskTag ScheduledTaskTag, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return scheduledTaskTag, errors.New("ScheduledTaskTagMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	re := regexp.MustCompile(scheduledTaskTagRegex)
	if !re.MatchString(stringValue) {
		return scheduledTaskTag, errors.New("InvalidScheduledTaskTag")
	}

	return ScheduledTaskTag(stringValue), nil
}

func (vo ScheduledTaskTag) String() string {
	return string(vo)
}
