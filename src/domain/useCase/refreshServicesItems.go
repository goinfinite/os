package useCase

import (
	"log/slog"

	"github.com/goinfinite/os/src/domain/repository"
)

func RefreshServicesItems(servicesCmdRepo repository.ServicesCmdRepo) {
	err := servicesCmdRepo.RefreshItems()
	if err != nil {
		slog.Error("RefreshServicesItemsInfraError", slog.Any("error", err))
	}
}
