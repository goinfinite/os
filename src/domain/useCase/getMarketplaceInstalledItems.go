package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func GetMarketplaceInstalledItems(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
) ([]entity.MarketplaceInstalledItem, error) {
	marketplaceInstalledItems, err := marketplaceQueryRepo.GetInstalledItems()
	if err != nil {
		log.Printf("GetMkplaceInstalledItemsError: %s", err.Error())
		return nil, errors.New("GetMkplaceInstalledItemsInfraError")
	}

	return marketplaceInstalledItems, nil
}
