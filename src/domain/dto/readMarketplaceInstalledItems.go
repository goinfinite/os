package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadMarketplaceInstalledItemsRequest struct {
	Pagination                       Pagination                                `json:"pagination"`
	MarketplaceInstalledItemId       *valueObject.MarketplaceItemId            `json:"marketplaceInstalledItemId,omitempty"`
	MarketplaceInstalledItemHostname *valueObject.Fqdn                         `json:"marketplaceInstalledItemHostname,omitempty"`
	MarketplaceInstalledItemType     *valueObject.MarketplaceItemType          `json:"marketplaceInstalledItemType,omitempty"`
	MarketplaceInstalledItemUuid     *valueObject.MarketplaceInstalledItemUuid `json:"marketplaceInstalledItemUuid,omitempty"`
}

type ReadMarketplaceInstalledItemsResponse struct {
	Pagination                Pagination                        `json:"pagination"`
	MarketplaceInstalledItems []entity.MarketplaceInstalledItem `json:"items"`
}
