package valueObject

type MarketplaceCatalogItemDataField struct {
	Key          DataFieldKey    `json:"key"`
	DefaultValue *DataFieldValue `json:"defaultValue,omitempty"`
	IsRequired   bool            `json:"isRequired"`
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
