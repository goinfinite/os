package useCase

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/repository"
)

func GetServices(
	servicesQueryRepo repository.ServicesQueryRepo,
) ([]entity.Service, error) {
	return servicesQueryRepo.Get()
}
