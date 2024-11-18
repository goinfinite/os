package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type MarketplaceQueryRepo interface {
	ReadCatalogItems(
		readDto dto.ReadMarketplaceCatalogItemsRequest,
	) (dto.ReadMarketplaceCatalogItemsResponse, error)
	ReadOneCatalogItem(
		readDto dto.ReadMarketplaceCatalogItemsRequest,
	) (entity.MarketplaceCatalogItem, error)
	ReadInstalledItems(
		readDto dto.ReadMarketplaceInstalledItemsRequest,
	) (dto.ReadMarketplaceInstalledItemsResponse, error)
	ReadOneInstalledItem(
		readDto dto.ReadMarketplaceInstalledItemsRequest,
	) (entity.MarketplaceInstalledItem, error)
}
