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

	mit := MktplaceItemType(value)
	if !mit.isValid() {
		return "", errors.New("InvalidMarketplaceItemType")
	}

	return MktplaceItemType(value), nil
}

func NewMktplaceItemTypePanic(value string) MktplaceItemType {
	mit, err := NewMktplaceItemType(value)
	if err != nil {
		panic(err)
	}

	return mit
}

func (mit MktplaceItemType) isValid() bool {
	return slices.Contains(ValidMktplaceItemTypes, string(mit))
}

func (mit MktplaceItemType) String() string {
	return string(mit)
}

func (mitPtr *MktplaceItemType) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	unquotedValue := strings.Trim(valueStr, "\"")

	mit, err := NewMktplaceItemType(unquotedValue)
	if err != nil {
		return err
	}

	*mitPtr = mit
	return nil
}
