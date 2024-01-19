package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func isPhpInstalled(
	runtimeQueryRepo repository.RuntimeQueryRepo,
	phpVersion valueObject.PhpVersion,
) bool {
	phpVersions, err := runtimeQueryRepo.GetPhpVersionsInstalled()
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
	updatePhpConfigsDto dto.UpdatePhpConfigs,
) error {
	isPhpInstalled := isPhpInstalled(
		runtimeQueryRepo,
		updatePhpConfigsDto.PhpVersion,
	)
	if !isPhpInstalled {
		return errors.New("PhpVersionNotInstalled")
	}

	_, err := vhostQueryRepo.GetByHostname(
		updatePhpConfigsDto.Hostname,
	)
	if err != nil {
		log.Printf("HostnameNotFound: %s", err.Error())
		return errors.New("HostnameNotFound")
	}

	err = runtimeCmdRepo.UpdatePhpVersion(
		updatePhpConfigsDto.Hostname,
		updatePhpConfigsDto.PhpVersion,
	)
	if err != nil {
		return errors.New("UpdatePhpVersionError")
	}

	err = runtimeCmdRepo.UpdatePhpSettings(
		updatePhpConfigsDto.Hostname,
		updatePhpConfigsDto.PhpSettings,
	)
	if err != nil {
		return errors.New("UpdatePhpSettingsError")
	}

	err = runtimeCmdRepo.UpdatePhpModules(
		updatePhpConfigsDto.Hostname,
		updatePhpConfigsDto.PhpModules,
	)
	if err != nil {
		return errors.New("UpdatePhpModulesError")
	}

	return nil
}
