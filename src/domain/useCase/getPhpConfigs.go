package useCase

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func GetPhpConfigs(
	runtimeQueryRepo repository.RuntimeQueryRepo,
	hostname valueObject.Fqdn,
) (entity.PhpConfigs, error) {
	return runtimeQueryRepo.GetPhpConfigs(hostname)
}
