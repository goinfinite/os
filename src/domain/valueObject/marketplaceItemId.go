package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type MarketplaceItemId int64

func NewMarketplaceItemId(value interface{}) (MarketplaceItemId, error) {
	marketplaceItemUid, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidMarketplaceItemId")
	}

	mii := MarketplaceItemId(marketplaceItemUid)
	if !mii.isValid() {
		return 0, errors.New("InvalidMarketplaceItemId")
	}

	return mii, nil
}

func NewMarketplaceItemIdPanic(value interface{}) MarketplaceItemId {
	mii, err := NewMarketplaceItemId(value)
	if err != nil {
		panic(err)
	}

	return mii
}

func (mii MarketplaceItemId) isValid() bool {
	return mii >= 1 && mii <= 999999999999
}

func (mii MarketplaceItemId) Get() int64 {
	return int64(mii)
}

func (mii MarketplaceItemId) String() string {
	return strconv.FormatInt(int64(mii), 10)
}
