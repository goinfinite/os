package infraHelper

import (
	"regexp"
)

func GetAllRegexGroupMatches(input string, regexExpression string) []string {
	matchesValues := []string{}

	regex := regexp.MustCompile(regexExpression)
	matches := regex.FindAllStringSubmatch(input, -1)

	if len(matches) == 0 {
		return matchesValues
	}

	for _, match := range matches {
		matchesValues = append(matchesValues, match[1])
	}

	return matchesValues
}
