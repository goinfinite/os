package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
)

func DeleteVirtualHost(
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
	runtimeCmdRepo repository.RuntimeCmdRepo,
	servicesQueryRepo repository.ServicesQueryRepo,
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo,
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
		slog.Error("ReadVirtualHostEntityError", slog.String("err", err.Error()))
		return errors.New("ReadVirtualHostEntityError")
	}

	isPhpWebServerInstalled, err := servicesQueryRepo.IsInstalled(
		valueObject.ServiceNamePhpWebServer,
	)
	if err != nil {
		slog.Error(
			"ReadPhpWebServerInstallationError", slog.String("err", err.Error()),
		)
		return errors.New("ReadPhpWebServerInstallationError")
	}
	if isPhpWebServerInstalled {
		err = runtimeCmdRepo.DeletePhpVirtualHost(targetVirtualHost.Hostname)
		if err != nil {
			slog.Error("DeletePhpVirtualHostError", slog.String("err", err.Error()))
			return errors.New("DeletePhpVirtualHostError")
		}
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
