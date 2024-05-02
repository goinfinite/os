package dto

import "github.com/speedianet/os/src/domain/valueObject"

type DeleteMarketplaceInstalledItem struct {
	InstalledId             valueObject.MarketplaceInstalledItemId
	ShouldUninstallServices bool
	ShouldRemoveFiles       bool
}

func NewDeleteMarketplaceInstalledItem(
	installedId valueObject.MarketplaceInstalledItemId,
	shouldUninstallServices bool,
	shouldRemoveFiles bool,
) DeleteMarketplaceInstalledItem {
	return DeleteMarketplaceInstalledItem{
		InstalledId:             installedId,
		ShouldUninstallServices: shouldUninstallServices,
		ShouldRemoveFiles:       shouldRemoveFiles,
	}
}
