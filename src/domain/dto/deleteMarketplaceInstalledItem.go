package dto

import "github.com/speedianet/os/src/domain/valueObject"

type DeleteMarketplaceInstalledItem struct {
	InstalledId             valueObject.MarketplaceItemId
	ShouldUninstallServices bool
	ShouldRemoveFiles       bool
}

func NewDeleteMarketplaceInstalledItem(
	installedId valueObject.MarketplaceItemId,
	shouldUninstallServices bool,
	shouldRemoveFiles bool,
) DeleteMarketplaceInstalledItem {
	return DeleteMarketplaceInstalledItem{
		InstalledId:             installedId,
		ShouldUninstallServices: shouldUninstallServices,
		ShouldRemoveFiles:       shouldRemoveFiles,
	}
}
