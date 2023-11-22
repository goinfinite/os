package infraHelper

import (
	"regexp"
)

func GetRegexCapturingGroups(input string, regex string) []string {
	re := regexp.MustCompile(regex)
	matches := re.FindAllStringSubmatch(input, -1)

	capturingGroups := []string{}
	for _, match := range matches {
		if len(match) < 1 {
			continue
		}
		capturingGroups = append(capturingGroups, match[1:]...)
	}
	return capturingGroups
}
