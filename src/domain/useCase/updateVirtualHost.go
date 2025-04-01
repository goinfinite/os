package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func UpdateVirtualHost(
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateVirtualHost,
) error {
	_, err := vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
		Hostname: &updateDto.Hostname,
	})
	if err != nil {
		slog.Debug("ReadVirtualHostError", slog.String("err", err.Error()))
		return errors.New("VirtualHostNotFound")
	}

	err = vhostCmdRepo.Update(updateDto)
	if err != nil {
		slog.Error("UpdateVirtualHostError", slog.String("err", err.Error()))
		return errors.New("UpdateVirtualHostInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).UpdateVirtualHost(updateDto)

	return nil
}
