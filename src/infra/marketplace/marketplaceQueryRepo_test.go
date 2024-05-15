package marketplaceInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestVirtualHostQueryRepo(t *testing.T) {
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
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
