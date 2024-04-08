package valueObject

import "errors"

type MarketplaceItemDataField struct {
	Key          DataFieldKey    `json:"key"`
	Value        DataFieldValue  `json:"value"`
	IsRequired   bool            `json:"isRequired"`
	DefaultValue *DataFieldValue `json:"defaultValue,omitempty"`
}

func NewMarketplaceItemDataField(
	key DataFieldKey,
	value DataFieldValue,
	isRequired bool,
	defaultValue *DataFieldValue,
) (MarketplaceItemDataField, error) {
	var marketplaceItemDataField MarketplaceItemDataField

	missingRequiredDefaultValue := !isRequired && defaultValue == nil
	if missingRequiredDefaultValue {
		return marketplaceItemDataField, errors.New("MissingRequiredDefaultValue")
	}

	return MarketplaceItemDataField{
		Key:          key,
		Value:        value,
		IsRequired:   isRequired,
		DefaultValue: defaultValue,
	}, nil
}

func NewMarketplaceItemDataFieldPanic(
	key DataFieldKey,
	value DataFieldValue,
	isRequired bool,
	defaultValue *DataFieldValue,
) MarketplaceItemDataField {
	marketplaceItemDataField, err := NewMarketplaceItemDataField(
		key,
		value,
		isRequired,
		defaultValue,
	)
	if err != nil {
		panic(err)
	}

	return marketplaceItemDataField
}
