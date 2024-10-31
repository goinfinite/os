package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadMarketplaceInstalledItemsRequest struct {
	Pagination Pagination                       `json:"pagination"`
	ItemId     *valueObject.MarketplaceItemId   `json:"itemId,omitempty"`
	ItemType   *valueObject.MarketplaceItemType `json:"itemType,omitempty"`
}

type ReadMarketplaceInstalledItemsResponse struct {
	Pagination Pagination                        `json:"pagination"`
	Items      []entity.MarketplaceInstalledItem `json:"items"`
}
