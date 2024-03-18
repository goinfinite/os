package mktplaceInfra

import "github.com/speedianet/os/src/domain/entity"

type MarketplaceQueryRepo struct{}

func (mktplace MarketplaceQueryRepo) Get() ([]entity.MarketplaceCatalog, error) {
	marketplaceCatalogList := []entity.MarketplaceCatalog{}

	return marketplaceCatalogList, nil
}
