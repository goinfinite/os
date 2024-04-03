package repository

import (
	"github.com/speedianet/os/src/domain/dto"
)

type MktplaceCatalogCmdRepo interface {
	InstallItem(installMktplaceCatalogItem dto.InstallMarketplaceCatalogItem) error
}
