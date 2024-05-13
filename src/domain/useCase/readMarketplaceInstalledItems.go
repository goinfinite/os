package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func ReadMarketplaceInstalledItems(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
) ([]entity.MarketplaceInstalledItem, error) {
	installedItems, err := marketplaceQueryRepo.ReadInstalledItems()
	if err != nil {
		log.Printf("ReadMarketplaceInstalledItemsError: %s", err.Error())
		return nil, errors.New("ReadMarketplaceInstalledItemsInfraError")
	}

	return installedItems, nil
}
