package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadMarketplaceCatalogItemsRequest struct {
	Pagination Pagination                       `json:"pagination"`
	ItemId     *valueObject.MarketplaceItemId   `json:"itemId,omitempty"`
	ItemSlug   *valueObject.MarketplaceItemSlug `json:"itemSlug,omitempty"`
	ItemName   *valueObject.MarketplaceItemName `json:"itemName,omitempty"`
	ItemType   *valueObject.MarketplaceItemType `json:"itemType,omitempty"`
}

type ReadMarketplaceCatalogItemsResponse struct {
	Pagination Pagination                      `json:"pagination"`
	Items      []entity.MarketplaceCatalogItem `json:"items"`
}
