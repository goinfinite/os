package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
)

type InstallMarketplaceCatalogItem struct {
	Hostname   valueObject.Fqdn                                  `json:"hostname"`
	Id         *valueObject.MarketplaceItemId                    `json:"id"`
	Slug       *valueObject.MarketplaceItemSlug                  `json:"slug"`
	UrlPath    *valueObject.UrlPath                              `json:"urlPath"`
	DataFields []valueObject.MarketplaceInstallableItemDataField `json:"dataFields"`
}

func NewInstallMarketplaceCatalogItem(
	hostname valueObject.Fqdn,
	id *valueObject.MarketplaceItemId,
	slug *valueObject.MarketplaceItemSlug,
	urlPath *valueObject.UrlPath,
	dataFields []valueObject.MarketplaceInstallableItemDataField,
) InstallMarketplaceCatalogItem {
	return InstallMarketplaceCatalogItem{
		Id:         id,
		Slug:       slug,
		Hostname:   hostname,
		UrlPath:    urlPath,
		DataFields: dataFields,
	}
}
