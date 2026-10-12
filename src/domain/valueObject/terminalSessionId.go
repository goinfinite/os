package valueObject

import (
	"errors"
	"regexp"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type TerminalSessionId string

var terminalSessionIdRegex = regexp.MustCompile(`^[a-f0-9]{16}$`)

func NewTerminalSessionId(value any) (
	terminalSessionId TerminalSessionId, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return terminalSessionId, errors.New("TerminalSessionIdMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !terminalSessionIdRegex.MatchString(stringValue) {
		return terminalSessionId, errors.New("InvalidTerminalSessionId")
	}

	return TerminalSessionId(stringValue), nil
}

func (vo TerminalSessionId) String() string {
	return string(vo)
}
