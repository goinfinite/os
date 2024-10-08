package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
)

type MarketplaceCmdRepo interface {
	InstallItem(installDto dto.InstallMarketplaceCatalogItem) error
	UninstallItem(deleteDto dto.DeleteMarketplaceInstalledItem) error
}
