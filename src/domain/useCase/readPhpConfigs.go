package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func ReadPhpConfigs(
	runtimeQueryRepo repository.RuntimeQueryRepo,
	hostname valueObject.Fqdn,
) (entity.PhpConfigs, error) {
	phpConfigs, err := runtimeQueryRepo.ReadPhpConfigs(hostname)
	if err != nil {
		log.Printf("ReadPhpConfigsError: %s", err.Error())
		return entity.PhpConfigs{}, errors.New("ReadPhpConfigsInfraError")
	}

	return phpConfigs, nil
}
