package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type VirtualHostType string

const (
	VirtualHostTypeTopLevel  VirtualHostType = "top-level"
	VirtualHostTypeSubdomain VirtualHostType = "subdomain"
	VirtualHostTypeAlias     VirtualHostType = "alias"
	VirtualHostTypeWildcard  VirtualHostType = "wildcard"
	VirtualHostTypePrimary   VirtualHostType = "primary"
)

var AvailableVirtualHostsTypes = []string{
	VirtualHostTypeTopLevel.String(), VirtualHostTypeSubdomain.String(),
	VirtualHostTypeAlias.String(),
}

func NewVirtualHostType(value interface{}) (vhostType VirtualHostType, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return vhostType, errors.New("VirtualHostTypeMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	stringValueVo := VirtualHostType(stringValue)
	switch stringValueVo {
	case VirtualHostTypeTopLevel, VirtualHostTypeSubdomain,
		VirtualHostTypeAlias, VirtualHostTypeWildcard:
		return stringValueVo, nil
	case VirtualHostTypePrimary:
		return VirtualHostTypeTopLevel, nil
	default:
		return vhostType, errors.New("InvalidVirtualHostType")
	}
}

func (vo VirtualHostType) String() string {
	return string(vo)
}
