package repository

import (
	"github.com/speedianet/os/src/domain/dto"
)

type MarketplaceCmdRepo interface {
	InstallItem(installDto dto.InstallMarketplaceCatalogItem) error
}
