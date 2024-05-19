package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type MarketplaceItemId uint

func NewMarketplaceItemId(value interface{}) (MarketplaceItemId, error) {
	marketplaceItemUid, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidMarketplaceItemId")
	}

	return MarketplaceItemId(marketplaceItemUid), nil
}

func NewMarketplaceItemIdPanic(value interface{}) MarketplaceItemId {
	vo, err := NewMarketplaceItemId(value)
	if err != nil {
		panic(err)
	}

	return vo
}

func (vo MarketplaceItemId) Get() uint {
	return uint(vo)
}

func (vo MarketplaceItemId) String() string {
	return strconv.FormatUint(uint(vo), 10)
}
