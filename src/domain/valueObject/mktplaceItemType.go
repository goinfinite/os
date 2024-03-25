package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type MktplaceItemType string

var ValidMktplaceItemTypes = []string{
	"app",
	"framework",
	"stack",
}

func NewMktplaceItemType(value string) (MktplaceItemType, error) {
	value = strings.ToLower(value)

	if !slices.Contains(ValidMktplaceItemTypes, value) {
		return "", errors.New("InvalidMarketplaceItemType")
	}

	return MktplaceItemType(value), nil
}

func NewMktplaceItemTypePanic(value string) MktplaceItemType {
	sn, err := NewMktplaceItemType(value)
	if err != nil {
		panic(err)
	}

	return sn
}

func (sn MktplaceItemType) String() string {
	return string(sn)
}
