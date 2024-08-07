package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func ReadServicesWithMetrics(
	servicesQueryRepo repository.ServicesQueryRepo,
) ([]dto.InstalledServiceWithMetrics, error) {
	servicesWithMetrics, err := servicesQueryRepo.ReadWithMetrics()
	if err != nil {
		slog.Info("ReadServicesWithMetricsError", slog.Any("err", err))
		return servicesWithMetrics, errors.New("ReadServicesWithMetricsInfraError")
	}

	return servicesWithMetrics, nil
}
