package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func GetMarketplaceCatalog(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
) ([]entity.MarketplaceCatalogItem, error) {
	marketplaceCatalogItems, err := marketplaceQueryRepo.GetItems()
	if err != nil {
		log.Printf("GetMkplaceCatalogItemsError: %s", err.Error())
		return nil, errors.New("GetMkplaceCatalogItemsInfraError")
	}

	return marketplaceCatalogItems, nil
}
