package voHelper

import "regexp"

func FindNamedGroupsMatches(regexPattern string, input string) map[string]string {
	namedGroupsMatchesMap := map[string]string{}
	compiledRegex := regexp.MustCompile(regexPattern)

	matchedGroups := compiledRegex.FindStringSubmatch(input)
	if len(matchedGroups) == 0 {
		return namedGroupsMatchesMap
	}
	matchedGroups = matchedGroups[1:]

	namedGroups := compiledRegex.SubexpNames()
	if len(namedGroups) == 0 {
		return namedGroupsMatchesMap
	}
	namedGroups = namedGroups[1:]

	for index, groupName := range namedGroups {
		if groupName != "" {
			namedGroupsMatchesMap[groupName] = matchedGroups[index]
		}
	}

	return namedGroupsMatchesMap
}
