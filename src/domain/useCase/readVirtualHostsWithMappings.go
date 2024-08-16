package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func ReadVirtualHostsWithMappings(
	mappingQueryRepo repository.MappingQueryRepo,
) ([]dto.VirtualHostWithMappings, error) {
	vhostsWithMappings, err := mappingQueryRepo.ReadWithMappings()
	if err != nil {
		slog.Error("ReadWithMappingsError", slog.Any("err", err))
		return vhostsWithMappings, errors.New("ReadWithMappingsInfraError")
	}

	return vhostsWithMappings, nil
}
