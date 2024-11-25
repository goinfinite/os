package valueObject

import (
	"errors"
	"log/slog"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const systemResourceIdentifierRegex string = `^sri://(?P<accountId>[\d]{1,64}):(?P<resourceType>[\w\_\-]{2,64})\/(?P<resourceId>[\w\_\.\-/]{1,256}|\*)$`

type SystemResourceIdentifier string

func NewSystemResourceIdentifier(
	value interface{},
) (sri SystemResourceIdentifier, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return sri, errors.New("SystemResourceIdentifierMustBeString")
	}

	re := regexp.MustCompile(systemResourceIdentifierRegex)
	if !re.MatchString(stringValue) {
		return sri, errors.New("InvalidSystemResourceIdentifier")
	}

	return SystemResourceIdentifier(stringValue), nil
}

// Note: this panic is solely to warn about the misuse of the VO, specifically for developer
// auditing, and has nothing to do with user input. This is not a standard and should not be
// followed for the development of other VOs.
func NewSystemResourceIdentifierIgnoreError(value interface{}) SystemResourceIdentifier {
	sri, err := NewSystemResourceIdentifier(value)
	if err != nil {
		panicMessage := "UnexpectedSystemResourceIdentifierCreationError"
		slog.Debug(panicMessage, slog.Any("value", value), slog.Any("error", err))
		panic(panicMessage)
	}

	return sri
}

func NewAccountSri(accountId AccountId) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://0:account/" + accountId.String(),
	)
}

func NewSecureAccessKeySri(
	accountId AccountId,
	SecureAccessKeyId SecureAccessKeyId,
) SystemResourceIdentifier {
	return NewSystemResourceIdentifierIgnoreError(
		"sri://" + accountId.String() + ":secureAccessKey/" +
			SecureAccessKeyId.String(),
	)
}

func (vo SystemResourceIdentifier) String() string {
	return string(vo)
}
