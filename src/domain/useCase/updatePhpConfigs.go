package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func isPhpVersionInstalled(
	runtimeQueryRepo repository.RuntimeQueryRepo,
	phpVersion valueObject.PhpVersion,
) bool {
	phpVersions, err := runtimeQueryRepo.ReadPhpVersionsInstalled()
	if err != nil {
		return false
	}

	for _, version := range phpVersions {
		if version == phpVersion {
			return true
		}
	}

	return false
}

func UpdatePhpConfigs(
	runtimeQueryRepo repository.RuntimeQueryRepo,
	runtimeCmdRepo repository.RuntimeCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	activityCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdatePhpConfigs,
) error {
	isPhpVersionInstalled := isPhpVersionInstalled(runtimeQueryRepo, updateDto.PhpVersion)
	if !isPhpVersionInstalled {
		return errors.New("PhpVersionNotInstalled")
	}

	_, err := vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
		Hostname: &updateDto.Hostname,
	})
	if err != nil {
		slog.Error("VirtualHostNotFound", slog.String("err", err.Error()))
		return errors.New("VirtualHostNotFound")
	}

	err = runtimeCmdRepo.UpdatePhpVersion(updateDto.Hostname, updateDto.PhpVersion)
	if err != nil {
		slog.Error("UpdatePhpVersionError", slog.String("err", err.Error()))
		return errors.New("UpdatePhpVersionInfraError")
	}
	securityActivityRecord := NewCreateSecurityActivityRecord(activityCmdRepo)
	securityActivityRecord.UpdatePhpConfigs(updateDto, "version")

	if len(updateDto.PhpModules) > 0 {
		err = runtimeCmdRepo.UpdatePhpModules(updateDto.Hostname, updateDto.PhpModules)
		if err != nil {
			slog.Error("UpdatePhpModulesError", slog.String("err", err.Error()))
			return errors.New("UpdatePhpModulesInfraError")
		}
		securityActivityRecord.UpdatePhpConfigs(updateDto, "modules")
	}

	if len(updateDto.PhpSettings) > 0 {
		err = runtimeCmdRepo.UpdatePhpSettings(updateDto.Hostname, updateDto.PhpSettings)
		if err != nil {
			slog.Error("UpdatePhpSettingsError", slog.String("err", err.Error()))
			return errors.New("UpdatePhpSettingsInfraError")
		}
		securityActivityRecord.UpdatePhpConfigs(updateDto, "settings")
	}

	return nil
}
