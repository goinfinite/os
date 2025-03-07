package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadVirtualHostsWithMappings(
	mappingQueryRepo repository.MappingQueryRepo,
) ([]dto.VirtualHostWithMappings, error) {
	vhostsWithMappings, err := mappingQueryRepo.ReadWithMappings()
	if err != nil {
		slog.Error("ReadWithMappingsError", slog.String("err", err.Error()))
		return vhostsWithMappings, errors.New("ReadWithMappingsInfraError")
	}

	return vhostsWithMappings, nil
}
