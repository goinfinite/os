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

	_, err := vhostQueryRepo.ReadByHostname(
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
		log.Printf("UpdatePhpVersionError: %s", err.Error())
		return errors.New("UpdatePhpVersionInfraError")
	}

	err = runtimeCmdRepo.UpdatePhpSettings(
		updatePhpConfigsDto.Hostname,
		updatePhpConfigsDto.PhpSettings,
	)
	if err != nil {
		log.Printf("UpdatePhpSettingsError: %s", err.Error())
		return errors.New("UpdatePhpSettingsInfraError")
	}

	err = runtimeCmdRepo.UpdatePhpModules(
		updatePhpConfigsDto.Hostname,
		updatePhpConfigsDto.PhpModules,
	)
	if err != nil {
		log.Printf("UpdatePhpModulesError: %s", err.Error())
		return errors.New("UpdatePhpModulesInfraError")
	}

	return nil
}
