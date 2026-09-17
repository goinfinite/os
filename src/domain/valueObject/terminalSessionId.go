package valueObject

import (
	"errors"
	"regexp"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type TerminalSessionId string

const terminalSessionIdRegexExpression = `^[a-f0-9]{16}$`

func NewTerminalSessionId(value interface{}) (
	terminalSessionId TerminalSessionId, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return terminalSessionId, errors.New("TerminalSessionIdMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	re := regexp.MustCompile(terminalSessionIdRegexExpression)
	if !re.MatchString(stringValue) {
		return terminalSessionId, errors.New("InvalidTerminalSessionId")
	}

	return TerminalSessionId(stringValue), nil
}

func (vo TerminalSessionId) String() string {
	return string(vo)
}
