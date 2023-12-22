package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type VirtualHostType string

var ValidVirtualHostTypes = []string{
	"top-level",
	"subdomain",
	"wildcard",
	"alias",
}

func NewVirtualHostType(value string) (VirtualHostType, error) {
	value = strings.ToLower(value)
	if !slices.Contains(ValidVirtualHostTypes, value) {
		return "", errors.New("InvalidVirtualHostType")
	}
	return VirtualHostType(value), nil
}

func NewVirtualHostTypePanic(value string) VirtualHostType {
	vt, err := NewVirtualHostType(value)
	if err != nil {
		panic(err)
	}
	return vt
}

func (vt VirtualHostType) String() string {
	return string(vt)
}
