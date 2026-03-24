package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
)

type ReadMarketplaceCatalogItemsRequest struct {
	Pagination                 tkDto.Pagination                 `json:"pagination"`
	MarketplaceCatalogItemId   *valueObject.MarketplaceItemId   `json:"marketplaceCatalogItemId,omitempty"`
	MarketplaceCatalogItemSlug *valueObject.MarketplaceItemSlug `json:"marketplaceCatalogItemSlug,omitempty"`
	MarketplaceCatalogItemName *valueObject.MarketplaceItemName `json:"marketplaceCatalogItemName,omitempty"`
	MarketplaceCatalogItemType *valueObject.MarketplaceItemType `json:"marketplaceCatalogItemType,omitempty"`
}

type ReadMarketplaceCatalogItemsResponse struct {
	Pagination              tkDto.Pagination                `json:"pagination"`
	MarketplaceCatalogItems []entity.MarketplaceCatalogItem `json:"marketplaceCatalogItems"`
}
