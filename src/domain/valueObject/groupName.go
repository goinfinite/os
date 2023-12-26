package valueObject

import (
	"errors"
	"regexp"
)

const groupNameRegexExpression = `^[a-zA-Z0-9_-]{1,32}$`

type GroupName string

func NewGroupName(value string) (GroupName, error) {
	groupName := GroupName(value)
	if !groupName.isValid() {
		return "", errors.New("InvalidGroupName")
	}

	return groupName, nil
}

func NewGroupNamePanic(value string) GroupName {
	groupName, err := NewGroupName(value)
	if err != nil {
		panic(err)
	}
	return groupName
}

func (groupName GroupName) isValid() bool {
	groupNameRegex := regexp.MustCompile(groupNameRegexExpression)
	return groupNameRegex.MatchString(string(groupName))
}

func (groupName GroupName) String() string {
	return string(groupName)
}
