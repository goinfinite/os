package repository

import (
	"github.com/speedianet/os/src/domain/dto"
)

type MarketplaceCmdRepo interface {
	InstallItem(installMarketplaceCatalogItem dto.InstallMarketplaceCatalogItem) error
}
