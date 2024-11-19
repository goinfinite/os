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

	t.Run("ReadCatalogItems", func(t *testing.T) {
		itemType, _ := valueObject.NewMarketplaceItemType("app")

		paginationDto := useCase.MarketplaceDefaultPagination
		sortBy, _ := valueObject.NewPaginationSortBy("id")
		sortDirection, _ := valueObject.NewPaginationSortDirection("desc")
		paginationDto.SortBy = &sortBy
		paginationDto.SortDirection = &sortDirection

		readCatalogItemRequestDto := dto.ReadMarketplaceCatalogItemsRequest{
			Pagination:                 paginationDto,
			MarketplaceCatalogItemType: &itemType,
		}

		responseDto, err := marketplaceQueryRepo.ReadCatalogItems(
			readCatalogItemRequestDto,
		)
		if err != nil {
			t.Errorf("ReadMarketplaceCatalogItemsError: %v", err)
			return
		}

		if len(responseDto.MarketplaceCatalogItems) == 0 {
			t.Error("NoCatalogItemsFound")
		}
	})

	t.Run("ReadOneCatalogItem", func(t *testing.T) {
		itemType, _ := valueObject.NewMarketplaceItemType("app")

		readCatalogItemRequestDto := dto.ReadMarketplaceCatalogItemsRequest{
			MarketplaceCatalogItemType: &itemType,
		}

		_, err := marketplaceQueryRepo.ReadOneCatalogItem(
			readCatalogItemRequestDto,
		)
		if err != nil {
			t.Errorf("ReadOneMarketplaceCatalogItemError: %v", err)
			return
		}
	})

	t.Run("ReadInstalledItems", func(t *testing.T) {
		itemType, _ := valueObject.NewMarketplaceItemType("app")

		paginationDto := useCase.MarketplaceDefaultPagination
		sortBy, _ := valueObject.NewPaginationSortBy("id")
		sortDirection, _ := valueObject.NewPaginationSortDirection("desc")
		paginationDto.SortBy = &sortBy
		paginationDto.SortDirection = &sortDirection

		readInstalledItemRequestDto := dto.ReadMarketplaceInstalledItemsRequest{
			Pagination:                   paginationDto,
			MarketplaceInstalledItemType: &itemType,
		}

		responseDto, err := marketplaceQueryRepo.ReadInstalledItems(
			readInstalledItemRequestDto,
		)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
			return
		}

		if len(responseDto.MarketplaceInstalledItems) == 0 {
			t.Error("NoInstalledItemsFound")
		}
	})

	t.Run("ReadOneInstalledItem", func(t *testing.T) {
		itemType, _ := valueObject.NewMarketplaceItemType("app")

		readInstalledItemRequestDto := dto.ReadMarketplaceInstalledItemsRequest{
			MarketplaceInstalledItemType: &itemType,
		}

		_, err := marketplaceQueryRepo.ReadOneInstalledItem(
			readInstalledItemRequestDto,
		)
		if err != nil {
			t.Errorf("ReadOneMarketplaceInstalledItemError: %v", err)
			return
		}
	})
}
