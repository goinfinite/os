package useCase

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
)

func InstallMarketplaceCatalogItem(
	mktplaceCatalogQueryRepo repository.MktplaceCatalogQueryRepo,
	mktplaceCatalogCmdRepo repository.MktplaceCatalogCmdRepo,
	vhostQueryRepo vhostInfra.VirtualHostQueryRepo,
	vhostCmdRepo vhostInfra.VirtualHostCmdRepo,
	installMktplaceCatalogItem dto.InstallMarketplaceCatalogItem,
) error {
	return nil
}
