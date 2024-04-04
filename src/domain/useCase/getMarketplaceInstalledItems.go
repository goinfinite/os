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
	installedItems, err := marketplaceQueryRepo.GetInstalledItems()
	if err != nil {
		log.Printf("GetMkplaceInstalledItemsError: %s", err.Error())
		return nil, errors.New("GetMkplaceInstalledItemsInfraError")
	}

	return installedItems, nil
}
