package valueObject

type MarketplaceCatalogItemDataField struct {
	Key          DataFieldKey     `json:"key"`
	Label        DataFieldLabel   `json:"label"`
	HtmlType     HtmlType         `json:"type"`
	DefaultValue *DataFieldValue  `json:"defaultValue,omitempty"`
	Options      []DataFieldValue `json:"options"`
	IsRequired   bool             `json:"isRequired"`
}

func NewMarketplaceCatalogItemDataField(
	key DataFieldKey,
	label DataFieldLabel,
	htmlType HtmlType,
	defaultValue *DataFieldValue,
	options []DataFieldValue,
	isRequired bool,
) (MarketplaceCatalogItemDataField, error) {
	return MarketplaceCatalogItemDataField{
		Key:          key,
		Label:        label,
		HtmlType:     htmlType,
		DefaultValue: defaultValue,
		Options:      options,
		IsRequired:   isRequired,
	}, nil
}
