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
	vHostReadResponse, err := vhostQueryRepo.Read(dto.ReadVirtualHostsRequest{
		Pagination: dto.PaginationUnpaginated,
	})
	if err != nil {
		slog.Error("ReadVirtualHostInfraError", slog.String("err", err.Error()))
		return errors.New("ReadVirtualHostInfraError")
	}

	existingVirtualHostHostnames := []valueObject.Fqdn{}
	for _, vhostEntity := range vHostReadResponse.VirtualHosts {
		if vhostEntity.Type == valueObject.VirtualHostTypeAlias {
			continue
		}

		if !slices.Contains(createDto.VirtualHostsHostnames, vhostEntity.Hostname) {
			continue
		}

		existingVirtualHostHostnames = append(existingVirtualHostHostnames, vhostEntity.Hostname)
	}

	if len(existingVirtualHostHostnames) == 0 {
		return errors.New("SpecifiedVirtualHostsNotFound")
	}

	createDto.VirtualHostsHostnames = existingVirtualHostHostnames

	sslPairId, err := sslCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateSslPairError", slog.String("err", err.Error()))
		return errors.New("CreateSslPairInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateSslPair(createDto, sslPairId)

	return nil
}
