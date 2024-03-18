package mktplaceInfra

import "github.com/speedianet/os/src/domain/entity"

type MarketplaceQueryRepo struct{}

func (mktplace MarketplaceQueryRepo) Get() ([]entity.MarketplaceItem, error) {
	marketplaceItems := []entity.MarketplaceItem{}

	return marketplaceItems, nil
}
