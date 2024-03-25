package mktplaceInfra

import "github.com/speedianet/os/src/domain/entity"

type MarketplaceQueryRepo struct{}

func (mktplace MarketplaceQueryRepo) Get() ([]entity.MarketplaceCatalogItem, error) {
	marketplaceCatalogItems := []entity.MarketplaceCatalogItem{}

	return marketplaceCatalogItems, nil
}
