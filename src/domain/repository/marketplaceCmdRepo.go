package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type MarketplaceCmdRepo interface {
	InstallItem(installDto dto.InstallMarketplaceCatalogItem) error
	UninstallItem(installedId valueObject.MarketplaceInstalledItemId) error
}
