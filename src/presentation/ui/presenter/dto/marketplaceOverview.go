package presenterDto

import (
	"github.com/goinfinite/os/src/domain/entity"
)

type CatalogItemsGroupedByType struct {
	Apps       []entity.MarketplaceCatalogItem
	Frameworks []entity.MarketplaceCatalogItem
	Stacks     []entity.MarketplaceCatalogItem
}

type MarketplaceOverview struct {
	ListType           string
	InstalledItemsList []entity.MarketplaceInstalledItem
	CatalogItemsList   CatalogItemsGroupedByType
}

func NewMarketplaceOverview(
	listType string,
	installedItemsList []entity.MarketplaceInstalledItem,
	catalogItemsList CatalogItemsGroupedByType,
) MarketplaceOverview {
	return MarketplaceOverview{
		ListType:           listType,
		InstalledItemsList: installedItemsList,
		CatalogItemsList:   catalogItemsList,
	}
}
