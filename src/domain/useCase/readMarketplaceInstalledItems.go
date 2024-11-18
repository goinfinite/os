package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadMarketplaceInstalledItems(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
	requestDto dto.ReadMarketplaceInstalledItemsRequest,
) (responseDto dto.ReadMarketplaceInstalledItemsResponse, err error) {
	responseDto, err = marketplaceQueryRepo.ReadInstalledItems(requestDto)
	if err != nil {
		slog.Error("ReadMarketplaceInstalledItemsError", slog.Any("error", err))
		return responseDto, errors.New("ReadMarketplaceInstalledItemsInfraError")
	}

	return responseDto, nil
}
