package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type ActivityRecordLevel string

var validActivityRecordLevels = []string{
	"DEBUG", "INFO", "WARN", "ERROR", "SEC",
}

func NewActivityRecordLevel(value interface{}) (
	activityRecordLevel ActivityRecordLevel, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return activityRecordLevel, errors.New("ActivityRecordLevelMustBeString")
	}
	stringValue = strings.ToUpper(stringValue)

	if !slices.Contains(validActivityRecordLevels, stringValue) {
		switch stringValue {
		case "SECURITY":
			stringValue = "SEC"
		case "WARNING":
			stringValue = "WARN"
		default:
			return activityRecordLevel, errors.New("InvalidActivityRecordLevel")
		}
	}

	return ActivityRecordLevel(stringValue), nil
}

func (vo ActivityRecordLevel) String() string {
	return string(vo)
}
