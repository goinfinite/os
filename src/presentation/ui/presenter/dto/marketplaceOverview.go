package presenterDto

import (
	"github.com/goinfinite/os/src/domain/entity"
	presenterValueObject "github.com/goinfinite/os/src/presentation/ui/presenter/valueObject"
)

type CatalogItemsGroupedByType struct {
	Apps       []entity.MarketplaceCatalogItem
	Frameworks []entity.MarketplaceCatalogItem
	Stacks     []entity.MarketplaceCatalogItem
}

type MarketplaceOverview struct {
	ListType           presenterValueObject.MarketplaceListType
	InstalledItemsList []entity.MarketplaceInstalledItem
	CatalogItemsList   CatalogItemsGroupedByType
}

func NewMarketplaceOverview(
	listType presenterValueObject.MarketplaceListType,
	installedItemsList []entity.MarketplaceInstalledItem,
	catalogItemsList CatalogItemsGroupedByType,
) MarketplaceOverview {
	return MarketplaceOverview{
		ListType:           listType,
		InstalledItemsList: installedItemsList,
		CatalogItemsList:   catalogItemsList,
	}
}
