package marketplaceInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
)

func TestVirtualHostQueryRepo(t *testing.T) {
	persistentDbSvc, _ := internalDbInfra.NewPersistentDatabaseService()
	marketplaceQueryRepo := NewMarketplaceQueryRepo(persistentDbSvc)
	testHelpers.LoadEnvVars()

	t.Run("Read", func(t *testing.T) {
		itemType, _ := valueObject.NewMarketplaceItemType("app")

		paginationDto := useCase.MarketplaceDefaultPagination
		sortBy, _ := valueObject.NewPaginationSortBy("id")
		sortDirection, _ := valueObject.NewPaginationSortDirection("desc")
		paginationDto.SortBy = &sortBy
		paginationDto.SortDirection = &sortDirection

		readDto := dto.ReadMarketplaceCatalogItemsRequest{
			Pagination: paginationDto,
			ItemType:   &itemType,
		}

		responseDto, err := marketplaceQueryRepo.ReadCatalogItems(readDto)
		if err != nil {
			t.Errorf("ReadMarketplaceItemsError: %v", err)
			return
		}

		if len(responseDto.Items) == 0 {
			t.Errorf("NoItemsFound")
		}
	})
}
