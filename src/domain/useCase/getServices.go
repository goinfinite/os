package useCase

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func GetServices(
	servicesQueryRepo repository.ServicesQueryRepo,
) ([]entity.Service, error) {
	return servicesQueryRepo.Get()
}
