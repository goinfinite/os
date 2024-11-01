package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadMarketplaceCatalogItemsRequest struct {
	Pagination Pagination                       `json:"pagination"`
	Id         *valueObject.MarketplaceItemId   `json:"id,omitempty"`
	Slug       *valueObject.MarketplaceItemSlug `json:"slug,omitempty"`
	Name       *valueObject.MarketplaceItemName `json:"name,omitempty"`
	Type       *valueObject.MarketplaceItemType `json:"type,omitempty"`
}

type ReadMarketplaceCatalogItemsResponse struct {
	Pagination Pagination                      `json:"pagination"`
	Items      []entity.MarketplaceCatalogItem `json:"items"`
}
