package valueObject

type MarketplaceCatalogItemDataField struct {
	Name         DataFieldName    `json:"name"`
	Label        DataFieldLabel   `json:"label"`
	Type         DataFieldType    `json:"type"`
	DefaultValue *DataFieldValue  `json:"defaultValue,omitempty"`
	Options      []DataFieldValue `json:"options,omitempty"`
	IsRequired   bool             `json:"isRequired"`
}

func NewMarketplaceCatalogItemDataField(
	name DataFieldName,
	label DataFieldLabel,
	fieldType DataFieldType,
	defaultValue *DataFieldValue,
	options []DataFieldValue,
	isRequired bool,
) (MarketplaceCatalogItemDataField, error) {
	return MarketplaceCatalogItemDataField{
		Name:         name,
		Label:        label,
		Type:         fieldType,
		DefaultValue: defaultValue,
		Options:      options,
		IsRequired:   isRequired,
	}, nil
}
