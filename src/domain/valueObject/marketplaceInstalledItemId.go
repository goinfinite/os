package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type MarketplaceInstalledItemId int64

func NewMarketplaceInstalledItemId(value interface{}) (MarketplaceInstalledItemId, error) {
	marketplaceItemUid, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidMarketplaceInstalledItemId")
	}

	miii := MarketplaceInstalledItemId(marketplaceItemUid)
	if !miii.isValid() {
		return 0, errors.New("InvalidMarketplaceInstalledItemId")
	}

	return miii, nil
}

func NewMarketplaceInstalledItemIdPanic(value interface{}) MarketplaceInstalledItemId {
	miii, err := NewMarketplaceInstalledItemId(value)
	if err != nil {
		panic(err)
	}

	return miii
}

func (miii MarketplaceInstalledItemId) isValid() bool {
	return miii >= 1 && miii <= 999999999999
}

func (miii MarketplaceInstalledItemId) Get() int64 {
	return int64(miii)
}

func (miii MarketplaceInstalledItemId) String() string {
	return strconv.FormatInt(int64(miii), 10)
}
