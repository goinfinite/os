package valueObject

type MarketplaceCatalogItemDataField struct {
	Key          DataFieldKey    `json:"key"`
	IsRequired   bool            `json:"isRequired"`
	DefaultValue *DataFieldValue `json:"defaultValue,omitempty"`
}

func NewMarketplaceCatalogItemDataField(
	key DataFieldKey,
	defaultValue *DataFieldValue,
	isRequired bool,
) (MarketplaceCatalogItemDataField, error) {
	return MarketplaceCatalogItemDataField{
		Key:          key,
		DefaultValue: defaultValue,
		IsRequired:   isRequired,
	}, nil
}
