package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type VirtualHostType string

var AvailableVirtualHostsTypes = []string{
	"top-level", "subdomain", "alias",
}
var validVirtualHostTypes = []string{
	"primary", "top-level", "subdomain", "wildcard", "alias",
}

func NewVirtualHostType(value interface{}) (vhostType VirtualHostType, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return vhostType, errors.New("VirtualHostTypeMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(validVirtualHostTypes, stringValue) {
		return vhostType, errors.New("InvalidVirtualHostType")
	}

	return VirtualHostType(stringValue), nil
}

func (vo VirtualHostType) String() string {
	return string(vo)
}
