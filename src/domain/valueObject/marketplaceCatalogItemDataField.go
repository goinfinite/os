package valueObject

import "errors"

type MarketplaceCatalogItemDataField struct {
	Key          DataFieldKey    `json:"key"`
	Value        DataFieldValue  `json:"value"`
	IsRequired   bool            `json:"isRequired"`
	DefaultValue *DataFieldValue `json:"defaultValue,omitempty"`
}

func NewMarketplaceCatalogItemDataField(
	key DataFieldKey,
	value DataFieldValue,
	isRequired bool,
	defaultValue *DataFieldValue,
) (MarketplaceCatalogItemDataField, error) {
	var marketplaceCatalogItemDataField MarketplaceCatalogItemDataField

	missingRequiredDefaultValue := !isRequired && defaultValue == nil
	if missingRequiredDefaultValue {
		return marketplaceCatalogItemDataField, errors.New("MissingRequiredDefaultValue")
	}

	return MarketplaceCatalogItemDataField{
		Key:          key,
		Value:        value,
		IsRequired:   isRequired,
		DefaultValue: defaultValue,
	}, nil
}

func NewMarketplaceCatalogItemDataFieldPanic(
	key DataFieldKey,
	value DataFieldValue,
	isRequired bool,
	defaultValue *DataFieldValue,
) MarketplaceCatalogItemDataField {
	marketplaceCatalogItemDataField, err := NewMarketplaceCatalogItemDataField(
		key,
		value,
		isRequired,
		defaultValue,
	)
	if err != nil {
		panic(err)
	}

	return marketplaceCatalogItemDataField
}
