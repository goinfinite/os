package valueObject

type MarketplaceInstalledItemDataField struct {
	Key   DataFieldKey   `json:"key"`
	Value DataFieldValue `json:"value"`
}

func NewMarketplaceInstalledItemDataField(
	key DataFieldKey,
	value DataFieldValue,
) (MarketplaceInstalledItemDataField, error) {
	return MarketplaceInstalledItemDataField{
		Key:   key,
		Value: value,
	}, nil
}

func NewMarketplaceInstalledItemDataFieldPanic(
	key DataFieldKey,
	value DataFieldValue,
) MarketplaceInstalledItemDataField {
	marketplaceInstalledItemDataField, err := NewMarketplaceInstalledItemDataField(
		key,
		value,
	)
	if err != nil {
		panic(err)
	}

	return marketplaceInstalledItemDataField
}
