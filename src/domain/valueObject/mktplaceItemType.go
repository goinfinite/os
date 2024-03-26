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
	mktplaceItemType, err := NewMktplaceItemType(value)
	if err != nil {
		panic(err)
	}

	return mktplaceItemType
}

func (mktplaceItemType MktplaceItemType) String() string {
	return string(mktplaceItemType)
}

func (mktplaceItemTypePtr *MktplaceItemType) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	unquotedValue := strings.Trim(valueStr, "\"")

	mktplaceItemType, err := NewMktplaceItemType(unquotedValue)
	if err != nil {
		return err
	}

	*mktplaceItemTypePtr = mktplaceItemType
	return nil
}
