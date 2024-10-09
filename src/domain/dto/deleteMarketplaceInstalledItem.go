package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteMarketplaceInstalledItem struct {
	InstalledId             valueObject.MarketplaceItemId
	ShouldUninstallServices bool
	ShouldRemoveFiles       bool
}

func NewDeleteMarketplaceInstalledItem(
	installedId valueObject.MarketplaceItemId,
	shouldUninstallServices bool,
) DeleteMarketplaceInstalledItem {
	return DeleteMarketplaceInstalledItem{
		InstalledId:             installedId,
		ShouldUninstallServices: shouldUninstallServices,
	}
}
