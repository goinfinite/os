package useCase

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func InstallMarketplaceCatalogItem(
	queryRepo repository.MktplaceCatalogQueryRepo,
	cmdRepo repository.MktplaceCatalogCmdRepo,
	dto dto.InstallMarketplaceCatalogItem,
) error {
	return nil
}
