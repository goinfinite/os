package cliController

import (
	"github.com/speedianet/os/src/domain/useCase"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	marketplaceInfra "github.com/speedianet/os/src/infra/marketplace"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

type MarketplaceController struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewMarketplaceController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MarketplaceController {
	return &MarketplaceController{
		persistentDbSvc: persistentDbSvc,
	}
}

func (controller MarketplaceController) GetCatalog() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "GetCatalogItems",
		Run: func(cmd *cobra.Command, args []string) {
			marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(
				controller.persistentDbSvc,
			)

			catalogItems, err := useCase.GetMarketplaceCatalog(marketplaceQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, catalogItems)
		},
	}
	return cmd
}
