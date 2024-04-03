package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func GetMarketplaceInstalledItems(
	mktplaceQueryRepo repository.MktplaceQueryRepo,
) ([]entity.MarketplaceInstalledItem, error) {
	mktplaceInstalledItems, err := mktplaceQueryRepo.GetInstalledItems()
	if err != nil {
		log.Printf("GetMkplaceInstalledItemsError: %s", err.Error())
		return nil, errors.New("GetMkplaceInstalledItemsInfraError")
	}

	return mktplaceInstalledItems, nil
}
