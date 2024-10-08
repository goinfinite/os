package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const relativeTimeRegex string = `^(?i)(\d+(?:\.\d+)?)\s*(second|minute|hour|day|week|month|year|s|m|h|d|w|M|y)(?:s?)\s*(ago|from now)?$`

type RelativeTime string

func NewRelativeTime(value interface{}) (relativeTime RelativeTime, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return relativeTime, errors.New("RelativeTimeMustBeString")
	}

	re := regexp.MustCompile(relativeTimeRegex)
	if !re.MatchString(stringValue) {
		return relativeTime, errors.New("InvalidRelativeTime")
	}

	return RelativeTime(stringValue), nil
}

func (vo RelativeTime) String() string {
	return string(vo)
}
