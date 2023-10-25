package infraHelper

import (
	"errors"
	"regexp"
)

func GetRegexFirstGroup(input string, regexExpression string) (string, error) {
	regex := regexp.MustCompile(regexExpression)
	match := regex.FindStringSubmatch(input)

	if len(match) < 2 {
		return "", errors.New("RegexGroupNotFound")
	}

	firstGroup := match[1]
	if len(firstGroup) == 0 {
		return "", errors.New("EmptyRegexGroup")
	}

	return firstGroup, nil
}
