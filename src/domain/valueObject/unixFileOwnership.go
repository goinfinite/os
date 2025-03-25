package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type UnixFileOwnership string

const UnixFileOwnershipAppWorkingDir = UnixFileOwnership("nobody:nogroup")

func NewUnixFileOwnership(
	value interface{},
) (fileOwnership UnixFileOwnership, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return fileOwnership, errors.New("UnixFileOwnershipValueMustBeString")
	}

	errorMessage := "InvalidUnixFileOwnership"

	stringValueParts := strings.Split(stringValue, ":")
	if len(stringValueParts) != 2 {
		return fileOwnership, errors.New(errorMessage)
	}

	fileOwner, err := NewUsername(stringValueParts[0])
	if err != nil {
		return fileOwnership, errors.New(errorMessage)
	}

	fileGroup, err := NewGroupName(stringValueParts[1])
	if err != nil {
		return fileOwnership, errors.New(errorMessage)
	}

	fileOwnershipStr := fileOwner.String() + ":" + fileGroup.String()
	return UnixFileOwnership(fileOwnershipStr), nil
}

func (vo UnixFileOwnership) String() string {
	return string(vo)
}
