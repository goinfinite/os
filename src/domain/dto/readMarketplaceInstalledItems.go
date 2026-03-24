package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type ReadMarketplaceInstalledItemsRequest struct {
	Pagination                       tkDto.Pagination                          `json:"pagination"`
	MarketplaceInstalledItemId       *valueObject.MarketplaceItemId            `json:"marketplaceInstalledItemId,omitempty"`
	MarketplaceInstalledItemHostname *tkValueObject.Fqdn                       `json:"marketplaceInstalledItemHostname,omitempty"`
	MarketplaceInstalledItemType     *valueObject.MarketplaceItemType          `json:"marketplaceInstalledItemType,omitempty"`
	MarketplaceInstalledItemUuid     *valueObject.MarketplaceInstalledItemUuid `json:"marketplaceInstalledItemUuid,omitempty"`
}

type ReadMarketplaceInstalledItemsResponse struct {
	Pagination                tkDto.Pagination                  `json:"pagination"`
	MarketplaceInstalledItems []entity.MarketplaceInstalledItem `json:"marketplaceInstalledItems"`
}
