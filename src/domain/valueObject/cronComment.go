package valueObject

import (
	"errors"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type CronComment string

func NewCronComment(value interface{}) (cronComment CronComment, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return cronComment, errors.New("CronCommentMustBeString")
	}

	if len(stringValue) < 2 {
		return cronComment, errors.New("CronCommentIsTooShort")
	}

	if len(stringValue) > 512 {
		return cronComment, errors.New("CronCommentIsTooLong")
	}

	return CronComment(stringValue), nil
}

func (vo CronComment) String() string {
	return string(vo)
}
