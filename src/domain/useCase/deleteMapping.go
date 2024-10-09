package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func DeleteMapping(
	queryRepo repository.MappingQueryRepo,
	cmdRepo repository.MappingCmdRepo,
	mappingId valueObject.MappingId,
) error {
	mapping, err := queryRepo.ReadById(mappingId)
	if err != nil {
		return errors.New("MappingNotFound")
	}

	err = cmdRepo.Delete(mappingId)
	if err != nil {
		log.Printf("DeleteMappingError: %v", err)
		return errors.New("DeleteMappingInfraError")
	}

	log.Printf(
		"Mapping '%s' from '%s' deleted.",
		mapping.Path.String(),
		mapping.Hostname.String(),
	)

	return nil
}
