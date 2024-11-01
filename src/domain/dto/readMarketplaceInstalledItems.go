package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadMarketplaceInstalledItemsRequest struct {
	Pagination       Pagination                                `json:"pagination"`
	Id               *valueObject.MarketplaceItemId            `json:"id,omitempty"`
	Hostname         *valueObject.Fqdn                         `json:"hostname,omitempty"`
	Type             *valueObject.MarketplaceItemType          `json:"type,omitempty"`
	InstallationUuid *valueObject.MarketplaceInstalledItemUuid `json:"installationUuid,omitempty"`
}

type ReadMarketplaceInstalledItemsResponse struct {
	Pagination Pagination                        `json:"pagination"`
	Items      []entity.MarketplaceInstalledItem `json:"items"`
}
