package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type InstallMarketplaceCatalogItem struct {
	Hostname          tkValueObject.Fqdn                                `json:"hostname"`
	Id                *valueObject.MarketplaceItemId                    `json:"id"`
	Slug              *valueObject.MarketplaceItemSlug                  `json:"slug"`
	UrlPath           *valueObject.UrlPath                              `json:"urlPath"`
	DataFields        []valueObject.MarketplaceInstallableItemDataField `json:"dataFields"`
	OperatorAccountId tkValueObject.AccountId                           `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress                           `json:"-"`
}

func NewInstallMarketplaceCatalogItem(
	hostname tkValueObject.Fqdn,
	id *valueObject.MarketplaceItemId,
	slug *valueObject.MarketplaceItemSlug,
	urlPath *valueObject.UrlPath,
	dataFields []valueObject.MarketplaceInstallableItemDataField,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
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
