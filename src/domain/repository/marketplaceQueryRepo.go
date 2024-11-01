package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type MarketplaceQueryRepo interface {
	ReadCatalogItems(
		readDto dto.ReadMarketplaceCatalogItemsRequest,
	) (dto.ReadMarketplaceCatalogItemsResponse, error)
	ReadUniqueCatalogItem(
		readDto dto.ReadMarketplaceCatalogItemsRequest,
	) (entity.MarketplaceCatalogItem, error)
	ReadInstalledItems(
		readDto dto.ReadMarketplaceInstalledItemsRequest,
	) (dto.ReadMarketplaceInstalledItemsResponse, error)
	ReadInstalledItemById(
		installedId valueObject.MarketplaceItemId,
	) (entity.MarketplaceInstalledItem, error)
}
