package valueObject

import (
	"gopkg.in/yaml.v3"
)

type MarketplaceItemDataField struct {
	Key   DataFieldKey   `json:"key"`
	Value DataFieldValue `json:"value"`
}

func NewMarketplaceItemDataField(
	key DataFieldKey,
	value DataFieldValue,
) MarketplaceItemDataField {
	return MarketplaceItemDataField{
		Key:   key,
		Value: value,
	}
}

func (dfPtr *MarketplaceItemDataField) UnmarshalYAML(value *yaml.Node) error {
	var valuesMap map[string]string
	err := value.Decode(&valuesMap)
	if err != nil {
		return err
	}

	dfKey, err := NewDataFieldKey(valuesMap["key"])
	if err != nil {
		return err
	}

	dfValue, err := NewDataFieldValue(valuesMap["value"])
	if err != nil {
		return err
	}

	*dfPtr = NewMarketplaceItemDataField(dfKey, dfValue)

	return nil
}
