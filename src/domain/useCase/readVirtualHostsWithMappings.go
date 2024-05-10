package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func ReadVirtualHostsWithMappings(
	mappingQueryRepo repository.MappingQueryRepo,
) ([]dto.VirtualHostWithMappings, error) {
	vhostsWithMappings, err := mappingQueryRepo.ReadWithMappings()
	if err != nil {
		log.Printf("ReadWithMappingsError: %s", err.Error())
		return vhostsWithMappings, errors.New("ReadWithMappingsInfraError")
	}

	return vhostsWithMappings, nil
}
