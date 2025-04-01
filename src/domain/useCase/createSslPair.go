package useCase

import (
	"errors"
	"log/slog"
	"slices"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func CreateSslPair(
	vhostQueryRepo repository.VirtualHostQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateSslPair,
) error {
	readVirtualHostsResponse, err := vhostQueryRepo.Read(dto.ReadVirtualHostsRequest{
		Pagination: dto.PaginationUnpaginated,
	})
	if err != nil {
		slog.Error("ReadVirtualHostInfraError", slog.String("err", err.Error()))
		return errors.New("ReadVirtualHostInfraError")
	}

	existingHostnames := []valueObject.Fqdn{}
	for _, vhostEntity := range readVirtualHostsResponse.VirtualHosts {
		if vhostEntity.Type == valueObject.VirtualHostTypeAlias {
			continue
		}

		if !slices.Contains(createDto.VirtualHostsHostnames, vhostEntity.Hostname) {
			continue
		}

		existingHostnames = append(existingHostnames, vhostEntity.Hostname)
	}

	if len(existingHostnames) == 0 {
		return errors.New("SpecifiedVirtualHostsNotFound")
	}

	createDto.VirtualHostsHostnames = existingHostnames

	sslPairId, err := sslCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateSslPairError", slog.String("err", err.Error()))
		return errors.New("CreateSslPairInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateSslPair(createDto, sslPairId)

	return nil
}
