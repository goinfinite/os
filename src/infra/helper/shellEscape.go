package infraHelper

import (
	"regexp"
	"strings"
	"unicode"
)

type ShellEscape struct {
}

func (helper ShellEscape) Quote(inputStr string) string {
	if len(inputStr) == 0 {
		return "''"
	}

	escapableCharsRegex := regexp.MustCompile(`[^\w@%+=:,./-]`)
	if !escapableCharsRegex.MatchString(inputStr) {
		return inputStr
	}

	return "'" + strings.ReplaceAll(inputStr, "'", "'\"'\"'") + "'"
}

func (helper ShellEscape) StripUnsafe(inputStr string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}

		return -1
	}, inputStr)
}
