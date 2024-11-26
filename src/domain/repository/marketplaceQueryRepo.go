package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type MarketplaceQueryRepo interface {
	ReadCatalogItems(
		dto.ReadMarketplaceCatalogItemsRequest,
	) (dto.ReadMarketplaceCatalogItemsResponse, error)
	ReadFirstCatalogItem(
		dto.ReadMarketplaceCatalogItemsRequest,
	) (entity.MarketplaceCatalogItem, error)
	ReadInstalledItems(
		dto.ReadMarketplaceInstalledItemsRequest,
	) (dto.ReadMarketplaceInstalledItemsResponse, error)
	ReadFirstInstalledItem(
		dto.ReadMarketplaceInstalledItemsRequest,
	) (entity.MarketplaceInstalledItem, error)
}
