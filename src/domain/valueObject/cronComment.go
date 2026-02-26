package valueObject

import (
	"errors"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type CronComment string

func NewCronComment(value interface{}) (cronComment CronComment, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return cronComment, errors.New("CronCommentMustBeString")
	}

	if len(stringValue) > 512 {
		return cronComment, errors.New("CronCommentIsTooLong")
	}

	return CronComment(stringValue), nil
}

func (vo CronComment) String() string {
	return string(vo)
}
