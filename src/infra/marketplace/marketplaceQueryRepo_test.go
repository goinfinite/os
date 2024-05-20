package marketplaceInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
)

func TestVirtualHostQueryRepo(t *testing.T) {
	persistentDbSvc, _ := internalDbInfra.NewPersistentDatabaseService()
	marketplaceQueryRepo := NewMarketplaceQueryRepo(persistentDbSvc)
	testHelpers.LoadEnvVars()

	t.Run("ReadCatalogItems", func(t *testing.T) {
		catalogItems, err := marketplaceQueryRepo.ReadCatalogItems()
		if err != nil {
			t.Errorf("ExpectingNoErrorButGot: %v", err)
		}

		if len(catalogItems) == 0 {
			t.Errorf("ExpectingEmptySliceButGot: %v", catalogItems)
		}
	})
}
