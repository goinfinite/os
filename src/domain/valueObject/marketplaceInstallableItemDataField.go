package valueObject

type MarketplaceInstallableItemDataField struct {
	Name  DataFieldName  `json:"name"`
	Value DataFieldValue `json:"value"`
}

func NewMarketplaceInstallableItemDataField(
	name DataFieldName, value DataFieldValue,
) MarketplaceInstallableItemDataField {
	return MarketplaceInstallableItemDataField{
		Name:  name,
		Value: value,
	}
}

func (vo MarketplaceInstallableItemDataField) String() string {
	return vo.Name.String() + ":" + vo.Value.String()
}
