package useCase

import (
	"errors"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
	"github.com/speedianet/sam/src/domain/valueObject"
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

func virtualHostExists(
	wsQueryRepo repository.WsQueryRepo,
	hostname valueObject.Fqdn,
) bool {
	hosts, err := wsQueryRepo.GetVirtualHosts()
	if err != nil {
		return false
	}

	for _, host := range hosts {
		if host == hostname {
			return true
		}
	}

	return false
}

func UpdatePhpConfigs(
	runtimeQueryRepo repository.RuntimeQueryRepo,
	runtimeCmdRepo repository.RuntimeCmdRepo,
	wsQueryRepo repository.WsQueryRepo,
	updatePhpConfigsDto dto.UpdatePhpConfigs,
) error {
	isPhpInstalled := isPhpInstalled(
		runtimeQueryRepo,
		updatePhpConfigsDto.PhpVersion,
	)
	if !isPhpInstalled {
		return errors.New("PhpVersionNotInstalled")
	}

	hostnameExists := virtualHostExists(
		wsQueryRepo,
		updatePhpConfigsDto.Hostname,
	)
	if !hostnameExists {
		return errors.New("HostnameNotFound")
	}

	err := runtimeCmdRepo.UpdatePhpVersion(
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
