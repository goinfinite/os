package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func ReadPhpConfigs(
	runtimeQueryRepo repository.RuntimeQueryRepo,
	hostname tkValueObject.Fqdn,
) (entity.PhpConfigs, error) {
	phpConfigs, err := runtimeQueryRepo.ReadPhpConfigs(hostname)
	if err != nil {
		log.Printf("ReadPhpConfigsError: %s", err.Error())
		return entity.PhpConfigs{}, errors.New("ReadPhpConfigsInfraError")
	}

	return phpConfigs, nil
}
