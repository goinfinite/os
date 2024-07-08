package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type UnixCommand string

func NewUnixCommand(value interface{}) (UnixCommand, error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("MarketplaceItemSlugValueMustBeString")
	}

	stringValue = strings.TrimSpace(stringValue)

	if len(stringValue) < 2 {
		return "", errors.New("UnixCommandTooShort")
	}

	if len(stringValue) > 4096 {
		return "", errors.New("UnixCommandTooLong")
	}

	return UnixCommand(stringValue), nil
}

func (vo UnixCommand) String() string {
	return string(vo)
}
