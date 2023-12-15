package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func DeleteMapping(
	queryRepo repository.VirtualHostQueryRepo,
	cmdRepo repository.VirtualHostCmdRepo,
	hostname valueObject.Fqdn,
	mappingId valueObject.MappingId,
) error {
	mappings, err := queryRepo.GetWithMappings()
	if err != nil {
		return errors.New("GetVirtualHostWithMappingsFailed")
	}

	var mapping *valueObject.Mapping
	for _, vhostWithMapping := range mappings {
		if vhostWithMapping.Hostname != hostname {
			continue
		}

		for _, vhostMapping := range vhostWithMapping.Mappings {
			if vhostMapping.Id != mappingId {
				continue
			}

			mapping = &vhostMapping
		}
	}
	if mapping == nil {
		return errors.New("MappingNotFound")
	}

	err = cmdRepo.DeleteMapping(hostname, *mapping)
	if err != nil {
		log.Printf("DeleteMappingError: %v", err)
		return errors.New("DeleteMappingInfraError")
	}

	log.Printf("Mapping '%v' deleted.", hostname)

	return nil
}
