package mktplaceInfra

import (
	"github.com/speedianet/os/src/domain/dto"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
)

type MktplaceCatalogCmdRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	queryRepo       *MktplaceCatalogQueryRepo
}

func NewMktplaceCatalogCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MktplaceCatalogCmdRepo {
	mktplaceCatalogQueryRepo := NewMktplaceCatalogQueryRepo(persistentDbSvc)

	return &MktplaceCatalogCmdRepo{
		persistentDbSvc: persistentDbSvc,
		queryRepo:       mktplaceCatalogQueryRepo,
	}
}

func (repo *MktplaceCatalogCmdRepo) InstallItem(
	installMktplaceCatalogItem dto.InstallMarketplaceCatalogItem,
) error {
	return nil
}
