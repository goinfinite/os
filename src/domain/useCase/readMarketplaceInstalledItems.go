package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadMarketplaceInstalledItems(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
	readDto dto.ReadMarketplaceInstalledItemsRequest,
) (dto.ReadMarketplaceInstalledItemsResponse, error) {
	responseDto, err := marketplaceQueryRepo.ReadInstalledItems(readDto)
	if err != nil {
		slog.Error("ReadMarketplaceInstalledItemsError", slog.Any("error", err))
		return responseDto, errors.New("ReadMarketplaceInstalledItemsInfraError")
	}

	return responseDto, nil
}
