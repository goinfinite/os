package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
)

type InstallMarketplaceCatalogItem struct {
	Hostname          valueObject.Fqdn                                  `json:"hostname"`
	Id                *valueObject.MarketplaceItemId                    `json:"id"`
	Slug              *valueObject.MarketplaceItemSlug                  `json:"slug"`
	UrlPath           *valueObject.UrlPath                              `json:"urlPath"`
	DataFields        []valueObject.MarketplaceInstallableItemDataField `json:"dataFields"`
	OperatorAccountId valueObject.AccountId                             `json:"-"`
	OperatorIpAddress valueObject.IpAddress                             `json:"-"`
}

func NewInstallMarketplaceCatalogItem(
	hostname valueObject.Fqdn,
	id *valueObject.MarketplaceItemId,
	slug *valueObject.MarketplaceItemSlug,
	urlPath *valueObject.UrlPath,
	dataFields []valueObject.MarketplaceInstallableItemDataField,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) InstallMarketplaceCatalogItem {
	return InstallMarketplaceCatalogItem{
		Id:                id,
		Slug:              slug,
		Hostname:          hostname,
		UrlPath:           urlPath,
		DataFields:        dataFields,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
