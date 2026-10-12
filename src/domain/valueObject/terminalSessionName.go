package valueObject

import (
	"errors"
	"regexp"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var terminalSessionNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9 ._()\-]{0,63}$`)

type TerminalSessionName string

func NewTerminalSessionName(value any) (
	terminalSessionName TerminalSessionName, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return terminalSessionName, errors.New("TerminalSessionNameMustBeString")
	}
	stringValue = strings.TrimSpace(stringValue)

	if !terminalSessionNameRegex.MatchString(stringValue) {
		return terminalSessionName, errors.New("InvalidTerminalSessionName")
	}

	return TerminalSessionName(stringValue), nil
}

func (vo TerminalSessionName) String() string {
	return string(vo)
}
