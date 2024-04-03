package repository

import (
	"github.com/speedianet/os/src/domain/dto"
)

type MktplaceCmdRepo interface {
	InstallItem(installMktplaceCatalogItem dto.InstallMarketplaceCatalogItem) error
}
