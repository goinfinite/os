package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadMarketplaceCatalogItemsRequest struct {
	Pagination                 Pagination                       `json:"pagination"`
	MarketplaceCatalogItemId   *valueObject.MarketplaceItemId   `json:"marketplaceCatalogItemId,omitempty"`
	MarketplaceCatalogItemSlug *valueObject.MarketplaceItemSlug `json:"MarketplaceCatalogItemSlug,omitempty"`
	MarketplaceCatalogItemName *valueObject.MarketplaceItemName `json:"marketplaceCatalogItemName,omitempty"`
	MarketplaceCatalogItemType *valueObject.MarketplaceItemType `json:"marketplaceCatalogItemType,omitempty"`
}

type ReadMarketplaceCatalogItemsResponse struct {
	Pagination              Pagination                      `json:"pagination"`
	MarketplaceCatalogItems []entity.MarketplaceCatalogItem `json:"marketplaceCatalogItems"`
}
