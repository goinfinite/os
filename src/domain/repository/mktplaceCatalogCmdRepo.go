package repository

import (
	"github.com/speedianet/os/src/domain/dto"
)

type MktplaceCatalogCmdRepo interface {
	Create(dto dto.InstallMarketplaceItem) error
}
