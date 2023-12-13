package voHelper

import "regexp"

func GetRegexNamedGroups(regex string, input string) map[string]string {
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(input)

	groupNames := re.SubexpNames()
	groupMap := make(map[string]string)
	for i, name := range groupNames {
		if i != 0 && name != "" {
			groupMap[name] = match[i]
		}
	}

	return groupMap
}
