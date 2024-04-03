package infraHelper

import (
	"errors"
	"regexp"
)

func GetRegexFirstGroup(input string, regexExpression string) ([]string, error) {
	matchesValues := []string{}

	regex := regexp.MustCompile(regexExpression)
	matches := regex.FindAllStringSubmatch(input, -1)

	if len(matches) == 0 {
		return matchesValues, errors.New("RegexGroupNotFound")
	}

	for _, match := range matches {
		matchesValues = append(matchesValues, match[1])
	}

	return matchesValues, nil
}
