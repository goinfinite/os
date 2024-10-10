package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const activityRecordRegex string = `^[A-Za-z]\w{2,128}$`

type ActivityRecordCode string

func NewActivityRecordCode(value interface{}) (code ActivityRecordCode, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return code, errors.New("ActivityRecordCodeMustBeString")
	}

	re := regexp.MustCompile(activityRecordRegex)
	if !re.MatchString(stringValue) {
		return code, errors.New("InvalidActivityRecordCode")
	}

	return ActivityRecordCode(stringValue), nil
}

func (vo ActivityRecordCode) String() string {
	return string(vo)
}
