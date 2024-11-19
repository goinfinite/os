package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type MarketplaceQueryRepo interface {
	ReadCatalogItems(
		requestDto dto.ReadMarketplaceCatalogItemsRequest,
	) (dto.ReadMarketplaceCatalogItemsResponse, error)
	ReadOneCatalogItem(
		requestDto dto.ReadMarketplaceCatalogItemsRequest,
	) (entity.MarketplaceCatalogItem, error)
	ReadInstalledItems(
		requestDto dto.ReadMarketplaceInstalledItemsRequest,
	) (dto.ReadMarketplaceInstalledItemsResponse, error)
	ReadOneInstalledItem(
		requestDto dto.ReadMarketplaceInstalledItemsRequest,
	) (entity.MarketplaceInstalledItem, error)
}
