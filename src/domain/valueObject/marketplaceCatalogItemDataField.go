package valueObject

type MarketplaceCatalogItemDataField struct {
	Name         DataFieldName          `json:"name"`
	Label        DataFieldLabel         `json:"label"`
	Type         DataFieldType          `json:"type"`
	SpecificType *DataFieldSpecificType `json:"specificType,omitempty"`
	DefaultValue *DataFieldValue        `json:"defaultValue,omitempty"`
	Options      []DataFieldValue       `json:"options,omitempty"`
	IsRequired   bool                   `json:"isRequired"`
}

func NewMarketplaceCatalogItemDataField(
	name DataFieldName,
	label DataFieldLabel,
	fieldType DataFieldType,
	fieldSpecificType *DataFieldSpecificType,
	defaultValue *DataFieldValue,
	options []DataFieldValue,
	isRequired bool,
) (MarketplaceCatalogItemDataField, error) {
	return MarketplaceCatalogItemDataField{
		Name:         name,
		Label:        label,
		Type:         fieldType,
		SpecificType: fieldSpecificType,
		DefaultValue: defaultValue,
		Options:      options,
		IsRequired:   isRequired,
	}, nil
}
