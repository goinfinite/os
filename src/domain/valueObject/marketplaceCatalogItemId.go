package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type MarketplaceCatalogItemId int64

func NewMarketplaceCatalogItemId(value interface{}) (MarketplaceCatalogItemId, error) {
	marketplaceItemUid, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidMarketplaceCatalogItemId")
	}

	mcii := MarketplaceCatalogItemId(marketplaceItemUid)
	if !mcii.isValid() {
		return 0, errors.New("InvalidMarketplaceCatalogItemId")
	}

	return mcii, nil
}

func NewMarketplaceCatalogItemIdPanic(value interface{}) MarketplaceCatalogItemId {
	mcii, err := NewMarketplaceCatalogItemId(value)
	if err != nil {
		panic(err)
	}

	return mcii
}

func (mcii MarketplaceCatalogItemId) isValid() bool {
	return mcii >= 1 && mcii <= 999999999999
}

func (mcii MarketplaceCatalogItemId) Get() int64 {
	return int64(mcii)
}

func (mcii MarketplaceCatalogItemId) String() string {
	return strconv.FormatInt(int64(mcii), 10)
}
