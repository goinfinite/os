package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func ReadServicesWithMetrics(
	servicesQueryRepo repository.ServicesQueryRepo,
) ([]dto.InstalledServiceWithMetrics, error) {
	servicesWithMetrics, err := servicesQueryRepo.ReadWithMetrics()
	if err != nil {
		log.Printf("ReadServicesWithMetricsError: %s", err.Error())
		return servicesWithMetrics, errors.New("ReadServicesWithMetricsInfraError")
	}

	return servicesWithMetrics, nil
}
