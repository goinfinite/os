package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

type MarketplaceItemType string

var ValidMarketplaceItemTypes = []string{
	"app",
	"framework",
	"stack",
}

func NewMarketplaceItemType(value string) (MarketplaceItemType, error) {
	value = strings.ToLower(value)

	mit := MarketplaceItemType(value)
	if !mit.isValid() {
		return "", errors.New("InvalidMarketplaceItemType")
	}

	return MarketplaceItemType(value), nil
}

func NewMarketplaceItemTypePanic(value string) MarketplaceItemType {
	mit, err := NewMarketplaceItemType(value)
	if err != nil {
		panic(err)
	}

	return mit
}

func (mit MarketplaceItemType) isValid() bool {
	return slices.Contains(ValidMarketplaceItemTypes, string(mit))
}

func (mit MarketplaceItemType) String() string {
	return string(mit)
}

func (mitPtr *MarketplaceItemType) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	unquotedValue := strings.Trim(valueStr, "\"")

	mit, err := NewMarketplaceItemType(unquotedValue)
	if err != nil {
		return err
	}

	*mitPtr = mit
	return nil
}

func (mitPtr *MarketplaceItemType) UnmarshalYAML(value *yaml.Node) error {
	var valueStr string
	err := value.Decode(&valueStr)
	if err != nil {
		return err
	}

	mit, err := NewMarketplaceItemType(valueStr)
	if err != nil {
		return err
	}

	*mitPtr = mit
	return nil
}
