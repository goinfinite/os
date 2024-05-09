package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func GetVirtualHostsWithMappings(
	mappingQueryRepo repository.MappingQueryRepo,
) ([]dto.VirtualHostWithMappings, error) {
	vhostsWithMappings, err := mappingQueryRepo.GetWithMappings()
	if err != nil {
		log.Printf("GetWithMappingsError: %s", err.Error())
		return vhostsWithMappings, errors.New("GetWithMappingsInfraError")
	}

	return vhostsWithMappings, nil
}
