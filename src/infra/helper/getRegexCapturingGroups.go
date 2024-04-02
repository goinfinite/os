package infraHelper

import (
	"regexp"
)

func GetRegexCapturingGroups(input string, regex string) map[string]string {
	namedGroupsWithValues := make(map[string]string)

	compiledRegex := regexp.MustCompile(regex)

	matches := compiledRegex.FindStringSubmatch(input)
	namedGroups := compiledRegex.SubexpNames()[1:]

	for namedGroupIndex, namedGroup := range namedGroups {
		namedGroupsWithValues[namedGroup] = matches[namedGroupIndex+1]
	}

	return namedGroupsWithValues
}
