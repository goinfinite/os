package infraHelper

import (
	"errors"
	"regexp"
)

func GetRegexUniqueMatch(input string, regexExpression string) (string, error) {
	regex := regexp.MustCompile(regexExpression)
	match := regex.FindStringSubmatch(input)

	if len(match) < 2 {
		return "", errors.New("StringMatchNotFound")
	}

	return match[1], nil
}
