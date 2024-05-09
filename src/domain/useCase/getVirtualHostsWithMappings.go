package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func GetVirtualHostsWithMappings(
	vhostQueryRepo repository.VirtualHostQueryRepo,
	mappingQueryRepo repository.MappingQueryRepo,
) ([]dto.VirtualHostWithMappings, error) {
	vhostsWithMappings := []dto.VirtualHostWithMappings{}

	vhosts, err := vhostQueryRepo.Get()
	if err != nil {
		log.Printf("GetVhostsError: %s", err.Error())
		return vhostsWithMappings, errors.New("GetVhostsInfraError")
	}

	for _, vhost := range vhosts {
		mappings, err := mappingQueryRepo.GetByHostname(vhost.Hostname)
		if err != nil {
			log.Printf("[%s] GetMappingsError: %s", vhost.Hostname, err.Error())
			continue
		}

		vhostsWithMappings = append(
			vhostsWithMappings,
			dto.NewVirtualHostWithMappings(
				vhost,
				mappings,
			),
		)
	}

	return vhostsWithMappings, nil
}
