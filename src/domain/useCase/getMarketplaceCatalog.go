package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func GetMarketplaceCatalog(
	mktplaceQueryRepo repository.MktplaceQueryRepo,
) ([]entity.MarketplaceCatalogItem, error) {
	mktplaceCatalogItems, err := mktplaceQueryRepo.GetItems()
	if err != nil {
		log.Printf("GetMkplaceCatalogItemsError: %s", err.Error())
		return nil, errors.New("GetMkplaceCatalogItemsInfraError")
	}

	return mktplaceCatalogItems, nil
}
