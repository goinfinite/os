package valueObject

type MarketplaceCatalogItemDataField struct {
	Name         DataFieldName    `json:"name"`
	Label        DataFieldLabel   `json:"label"`
	HtmlType     HtmlType         `json:"type"`
	DefaultValue *DataFieldValue  `json:"defaultValue,omitempty"`
	Options      []DataFieldValue `json:"options,omitempty"`
	IsRequired   bool             `json:"isRequired"`
}

func NewMarketplaceCatalogItemDataField(
	name DataFieldName,
	label DataFieldLabel,
	htmlType HtmlType,
	defaultValue *DataFieldValue,
	options []DataFieldValue,
	isRequired bool,
) (MarketplaceCatalogItemDataField, error) {
	return MarketplaceCatalogItemDataField{
		Name:         name,
		Label:        label,
		HtmlType:     htmlType,
		DefaultValue: defaultValue,
		Options:      options,
		IsRequired:   isRequired,
	}, nil
}
