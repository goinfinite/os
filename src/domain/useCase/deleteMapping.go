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
	vhostHostname valueObject.Fqdn,
	mappingId valueObject.MappingId,
) error {
	_, err := queryRepo.GetMappingById(vhostHostname, mappingId)
	if err != nil {
		return errors.New("MappingNotFound")
	}

	err = cmdRepo.DeleteMapping(vhostHostname, mappingId)
	if err != nil {
		log.Printf("DeleteMappingError: %v", err)
		return errors.New("DeleteMappingInfraError")
	}

	log.Printf("Mapping '%v' deleted.", vhostHostname)

	return nil
}
