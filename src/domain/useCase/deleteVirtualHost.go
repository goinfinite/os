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
	isPrimaryHostname := deleteDto.Hostname.String() == deleteDto.PrimaryVirtualHost.String()
	if isPrimaryHostname {
		return errors.New("PrimaryVirtualHostCannotBeDeleted")
	}

	vhost, err := vhostQueryRepo.ReadByHostname(deleteDto.Hostname)
	if err != nil {
		return errors.New("VirtualHostNotFound")
	}

	err = vhostCmdRepo.Delete(vhost)
	if err != nil {
		slog.Error("DeleteVirtualHostError", slog.Any("err", err))
		return errors.New("DeleteVirtualHostInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteVirtualHost(deleteDto)

	return nil
}
