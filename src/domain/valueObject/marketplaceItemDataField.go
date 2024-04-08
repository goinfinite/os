package valueObject

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
) MarketplaceItemDataField {
	return MarketplaceItemDataField{
		Key:          key,
		Value:        value,
		IsRequired:   isRequired,
		DefaultValue: defaultValue,
	}
}
