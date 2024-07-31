package valueObject

type MarketplaceInstallableItemDataField struct {
	Name  DataFieldName  `json:"name"`
	Value DataFieldValue `json:"value"`
}

func NewMarketplaceInstallableItemDataField(
	name DataFieldName, value DataFieldValue,
) (MarketplaceInstallableItemDataField, error) {
	return MarketplaceInstallableItemDataField{
		Name:  name,
		Value: value,
	}, nil
}

func (vo MarketplaceInstallableItemDataField) String() string {
	return vo.Name.String() + ":" + vo.Value.String()
}
