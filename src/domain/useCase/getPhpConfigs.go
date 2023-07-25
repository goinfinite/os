package useCase

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/repository"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func GetPhpConfigs(
	runtimeQueryRepo repository.RuntimeQueryRepo,
	hostname valueObject.Fqdn,
) (entity.PhpConfigs, error) {
	return runtimeQueryRepo.GetPhpConfigs(hostname)
}
