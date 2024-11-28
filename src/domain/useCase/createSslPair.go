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
	sslCmdRepo repository.SslCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateSslPair,
) error {
	existingVhosts, err := vhostQueryRepo.Read()
	if err != nil {
		slog.Error("ReadVhostsError", slog.Any("err", err))
		return errors.New("ReadVhostsInfraError")
	}

	if len(existingVhosts) == 0 {
		return errors.New("VhostsNotFound")
	}

	validSslVirtualHostsHostnames := []valueObject.Fqdn{}
	for _, vhost := range existingVhosts {
		if vhost.Type.String() == "alias" {
			continue
		}

		if slices.Contains(createDto.VirtualHostsHostnames, vhost.Hostname) {
			validSslVirtualHostsHostnames = append(
				validSslVirtualHostsHostnames, vhost.Hostname,
			)
		}
	}

	if len(validSslVirtualHostsHostnames) == 0 {
		return errors.New("VhostDoesNotExists")
	}

	createDto.VirtualHostsHostnames = validSslVirtualHostsHostnames

	createdSslPairId, err := sslCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateSslPairError", slog.Any("err", err))
		return errors.New("CreateSslPairInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateSslPair(createDto, createdSslPairId)

	return nil
}
