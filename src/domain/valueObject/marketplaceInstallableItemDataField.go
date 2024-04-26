package valueObject

type MarketplaceInstallableItemDataField struct {
	Key   DataFieldKey   `json:"key"`
	Value DataFieldValue `json:"value"`
}

func NewMarketplaceInstallableItemDataField(
	key DataFieldKey,
	value DataFieldValue,
) (MarketplaceInstallableItemDataField, error) {
	return MarketplaceInstallableItemDataField{
		Key:   key,
		Value: value,
	}, nil
}

func NewMarketplaceInstallableItemDataFieldPanic(
	key DataFieldKey,
	value DataFieldValue,
) MarketplaceInstallableItemDataField {
	marketplaceInstallableItemDataField, err := NewMarketplaceInstallableItemDataField(
		key,
		value,
	)
	if err != nil {
		panic(err)
	}

	return marketplaceInstallableItemDataField
}
