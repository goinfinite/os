package useCase

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
)

func InstallMarketplaceCatalogItem(
	queryRepo repository.MktplaceCatalogQueryRepo,
	cmdRepo repository.MktplaceCatalogCmdRepo,
	vhostQueryRepo vhostInfra.VirtualHostQueryRepo,
	dto dto.InstallMarketplaceCatalogItem,
) error {
	return nil
}
