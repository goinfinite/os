package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func GetPhpConfigs(
	runtimeQueryRepo repository.RuntimeQueryRepo,
	hostname valueObject.Fqdn,
) (entity.PhpConfigs, error) {
	phpConfigs, err := runtimeQueryRepo.GetPhpConfigs(hostname)
	if err != nil {
		log.Printf("GetPhpConfigsError: %s", err.Error())
		return entity.PhpConfigs{}, errors.New("GetPhpConfigsInfraError")
	}

	return phpConfigs, nil
}
