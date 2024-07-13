package valueObject

type MarketplaceInstallableItemDataField struct {
	Name  DataFieldName  `json:"name"`
	Value DataFieldValue `json:"value"`
}

func NewMarketplaceInstallableItemDataField(
	name DataFieldName,
	value DataFieldValue,
) (MarketplaceInstallableItemDataField, error) {
	return MarketplaceInstallableItemDataField{
		Name:  name,
		Value: value,
	}, nil
}

func NewMarketplaceInstallableItemDataFieldPanic(
	name DataFieldName,
	value DataFieldValue,
) MarketplaceInstallableItemDataField {
	marketplaceInstallableItemDataField, err := NewMarketplaceInstallableItemDataField(
		name,
		value,
	)
	if err != nil {
		panic(err)
	}

	return marketplaceInstallableItemDataField
}

func (vo MarketplaceInstallableItemDataField) String() string {
	return vo.Name.String() + ":" + vo.Value.String()
}
