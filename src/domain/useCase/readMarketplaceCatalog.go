package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func ReadMarketplaceCatalog(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
) ([]entity.MarketplaceCatalogItem, error) {
	catalogItems, err := marketplaceQueryRepo.ReadCatalogItems()
	if err != nil {
		log.Printf("ReadMarketplaceCatalogItemsError: %s", err.Error())
		return nil, errors.New("ReadMarketplaceCatalogItemsInfraError")
	}

	return catalogItems, nil
}
