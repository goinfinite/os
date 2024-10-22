package presenterDto

import (
	"github.com/goinfinite/os/src/domain/entity"
	presenterValueObject "github.com/goinfinite/os/src/presentation/ui/presenter/valueObject"
)

type MarketplaceOverview struct {
	VirtualHostsHostnames []string
	ListType              presenterValueObject.MarketplaceListType
	InstalledItemsList    []entity.MarketplaceInstalledItem
	CatalogItemsList      []entity.MarketplaceCatalogItem
}

func NewMarketplaceOverview(
	virtualHostsHostnames []string,
	listType presenterValueObject.MarketplaceListType,
	installedItemsList []entity.MarketplaceInstalledItem,
	catalogItemsList []entity.MarketplaceCatalogItem,
) MarketplaceOverview {
	return MarketplaceOverview{
		VirtualHostsHostnames: virtualHostsHostnames,
		ListType:              listType,
		InstalledItemsList:    installedItemsList,
		CatalogItemsList:      catalogItemsList,
	}
}
