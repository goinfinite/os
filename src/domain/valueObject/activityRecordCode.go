package valueObject

import (
	"errors"
	"slices"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type ActivityRecordCode string

var validActivityRecordCodes = []string{
	"LoginFailed", "LoginSuccessful",
	"AccountCreated", "AccountDeleted", "AccountPasswordUpdated",
	"AccountApiKeyUpdated", "AccountQuotaUpdated",
	"UnauthorizedAccess",
}

func NewActivityRecordCode(value interface{}) (
	activityRecordCode ActivityRecordCode, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return activityRecordCode, errors.New("ActivityRecordCodeMustBeString")
	}

	if !slices.Contains(validActivityRecordCodes, stringValue) {
		return activityRecordCode, errors.New("InvalidActivityRecordCode")
	}

	return ActivityRecordCode(stringValue), nil
}

func (vo ActivityRecordCode) String() string {
	return string(vo)
}
