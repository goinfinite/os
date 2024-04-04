package valueObject

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

type MarketplaceItemDescription string

func NewMarketplaceItemDescription(value string) (MarketplaceItemDescription, error) {
	if len(value) < 2 {
		return "", errors.New("MarketplaceItemDescriptionTooSmall")
	}

	if len(value) > 512 {
		return "", errors.New("MarketplaceItemDescriptionTooBig")
	}

	return MarketplaceItemDescription(value), nil
}

func NewMarketplaceItemDescriptionPanic(value string) MarketplaceItemDescription {
	mid, err := NewMarketplaceItemDescription(value)
	if err != nil {
		panic(err)
	}

	return mid
}

func (mid MarketplaceItemDescription) String() string {
	return string(mid)
}

func (midPtr *MarketplaceItemDescription) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	unquotedValue := strings.Trim(valueStr, "\"")

	mid, err := NewMarketplaceItemDescription(unquotedValue)
	if err != nil {
		return err
	}

	*midPtr = mid
	return nil
}

func (midPtr *MarketplaceItemDescription) UnmarshalYAML(value *yaml.Node) error {
	var valueStr string
	err := value.Decode(&valueStr)
	if err != nil {
		return err
	}

	mid, err := NewMarketplaceItemDescription(valueStr)
	if err != nil {
		return err
	}

	*midPtr = mid
	return nil
}
