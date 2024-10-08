package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const groupNameRegexExpression = `^[a-zA-Z0-9_-]{1,32}$`

type GroupName string

func NewGroupName(value interface{}) (groupName GroupName, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return groupName, errors.New("GroupNameMustBeString")
	}

	re := regexp.MustCompile(groupNameRegexExpression)
	if !re.MatchString(stringValue) {
		return groupName, errors.New("InvalidGroupName")
	}

	return GroupName(stringValue), nil
}

func (vo GroupName) String() string {
	return string(vo)
}
