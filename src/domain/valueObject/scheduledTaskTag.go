package valueObject

import (
	"errors"
	"regexp"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var scheduledTaskTagRegex = regexp.MustCompile(`^[a-zA-Z][\w\-]{1,256}$`)

type ScheduledTaskTag string

func NewScheduledTaskTag(value any) (
	scheduledTaskTag ScheduledTaskTag, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return scheduledTaskTag, errors.New("ScheduledTaskTagMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !scheduledTaskTagRegex.MatchString(stringValue) {
		return scheduledTaskTag, errors.New("InvalidScheduledTaskTag")
	}

	return ScheduledTaskTag(stringValue), nil
}

func (vo ScheduledTaskTag) String() string {
	return string(vo)
}
