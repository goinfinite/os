package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
)

type MarketplaceCmdRepo interface {
	InstallItem(dto.InstallMarketplaceCatalogItem) error
	UninstallItem(dto.DeleteMarketplaceInstalledItem) error
	RefreshCatalogItems() error
}
