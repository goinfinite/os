package mktplaceInfra

import (
	"github.com/speedianet/os/src/domain/dto"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
)

type MktplaceCatalogCmdRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewMktplaceCatalogCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MktplaceCatalogCmdRepo {
	return &MktplaceCatalogCmdRepo{
		persistentDbSvc: persistentDbSvc,
	}
}

func (repo *MktplaceCatalogCmdRepo) Create(
	dto dto.InstallMarketplaceItem,
) error {
	return nil
}
