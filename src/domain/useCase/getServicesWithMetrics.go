package useCase

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func GetServicesWithMetrics(
	servicesQueryRepo repository.ServicesQueryRepo,
) ([]dto.ServiceWithMetrics, error) {
	return servicesQueryRepo.GetWithMetrics()
}
