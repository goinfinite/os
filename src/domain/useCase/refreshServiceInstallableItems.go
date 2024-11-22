package useCase

import (
	"log/slog"

	"github.com/goinfinite/os/src/domain/repository"
)

const RefreshServiceInstallableItemsAmountPerDay int = 1

func RefreshServiceInstallableItems(servicesCmdRepo repository.ServicesCmdRepo) {
	err := servicesCmdRepo.RefreshInstallableItems()
	if err != nil {
		slog.Error("RefreshServiceInstallableItemsError", slog.Any("error", err))
	}
}
