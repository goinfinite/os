package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func DeleteMapping(
	queryRepo repository.MappingQueryRepo,
	cmdRepo repository.MappingCmdRepo,
	mappingId valueObject.MappingId,
) error {
	mapping, err := queryRepo.GetById(mappingId)
	if err != nil {
		return errors.New("MappingNotFound")
	}

	err = cmdRepo.DeleteMapping(mappingId)
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
