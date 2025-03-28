package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func DeleteVirtualHost(
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteVirtualHost,
) error {
	isPrimary := true
	primaryVirtualHost, err := vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
		IsPrimary: &isPrimary,
	})
	if err != nil {
		slog.Error("ReadPrimaryVirtualHostError", slog.String("err", err.Error()))
		return errors.New("ReadPrimaryVirtualHostError")
	}

	if primaryVirtualHost.Hostname == deleteDto.Hostname {
		return errors.New("PrimaryVirtualHostCannotBeDeleted")
	}

	targetVirtualHost, err := vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
		Hostname: &deleteDto.Hostname,
	})
	if err != nil {
		slog.Error("ReadVirtualHostError", slog.String("err", err.Error()))
		return errors.New("ReadVirtualHostError")
	}

	err = vhostCmdRepo.Delete(targetVirtualHost.Hostname)
	if err != nil {
		slog.Error("DeleteVirtualHostError", slog.String("err", err.Error()))
		return errors.New("DeleteVirtualHostInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteVirtualHost(deleteDto)

	return nil
}
