package runtimeInfra

import (
	"encoding/json"
	"errors"
	"log"
	"slices"
	"strings"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
)

type RuntimeQueryRepo struct {
}

func (repo RuntimeQueryRepo) GetVirtualHostPhpConfFilePath(
	hostname valueObject.Fqdn,
) (vhostPhpConfFilePath valueObject.UnixFilePath, err error) {
	primaryVhostPhpConfFilePathStr := "/app/conf/php-webserver/primary.conf"
	vhostPhpConfFilePathStr := "/app/conf/php-webserver/" + hostname.String() + ".conf"

	primaryVirtualHostHostname, err := infraHelper.ReadPrimaryVirtualHostHostname()
	if err != nil {
		return vhostPhpConfFilePath, errors.New("PrimaryVhostNotFound: " + err.Error())
	}

	if hostname == primaryVirtualHostHostname {
		vhostPhpConfFilePathStr = primaryVhostPhpConfFilePathStr
	}

	vhostPhpConfFilePath, err = valueObject.NewUnixFilePath(vhostPhpConfFilePathStr)
	if err != nil {
		return vhostPhpConfFilePath, err
	}

	if !infraHelper.FileExists(vhostPhpConfFilePathStr) {
		return vhostPhpConfFilePath, errors.New("VirtualHostNotFound")
	}

	return vhostPhpConfFilePath, nil
}

func (repo RuntimeQueryRepo) ReadPhpVersionsInstalled() (
	phpVersions []valueObject.PhpVersion, err error,
) {
	output, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "awk",
		Args: []string{
			"/extprocessor lsphp/{print $2}", infraEnvs.PhpWebserverMainConfFilePath,
		},
	})
	if err != nil {
		return phpVersions, errors.New("GetPhpVersionFromFileFailed: " + err.Error())
	}

	for _, version := range strings.Split(output, "\n") {
		if version == "" {
			continue
		}

		version = strings.Replace(version, "lsphp", "", 1)
		phpVersion, err := valueObject.NewPhpVersion(version)
		if err != nil {
			continue
		}

		phpVersions = append(phpVersions, phpVersion)
	}

	return phpVersions, nil
}

func (repo RuntimeQueryRepo) ReadPhpVersion(
	hostname valueObject.Fqdn,
) (phpVersion entity.PhpVersion, err error) {
	vhostPhpConfFilePath, err := repo.GetVirtualHostPhpConfFilePath(hostname)
	if err != nil {
		return phpVersion, err
	}

	currentPhpVersionStr, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "awk",
		Args: []string{
			"/lsapi:lsphp/ {gsub(/[^0-9]/, \"\", $2); print $2}",
			vhostPhpConfFilePath.String(),
		},
	})
	if err != nil {
		return phpVersion, errors.New("GetCurrentPhpVersionFromFileFailed: " + err.Error())
	}

	currentPhpVersion, err := valueObject.NewPhpVersion(currentPhpVersionStr)
	if err != nil {
		return phpVersion, errors.New("PhpVersionUnknown: " + err.Error())
	}

	phpVersions, err := repo.ReadPhpVersionsInstalled()
	if err != nil {
		return phpVersion, errors.New("GetPhpVersionsInstalledFailed: " + err.Error())
	}

	phpVersion = entity.NewPhpVersion(currentPhpVersion, phpVersions)
	return phpVersion, nil
}

func (repo RuntimeQueryRepo) getPhpTimezones() (timezones []string, err error) {
	timezonesRaw, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "php",
		Args:    []string{"-r", "echo json_encode(DateTimeZone::listIdentifiers());"},
	})
	if err != nil {
		return timezones, errors.New("GetPhpTimezonesFailed: " + err.Error())
	}

	err = json.Unmarshal([]byte(timezonesRaw), &timezones)
	if err != nil {
		return timezones, errors.New("ParsePhpTimezonesFailed: " + err.Error())
	}

	return timezones, nil
}

func (repo RuntimeQueryRepo) phpSettingFactory(
	setting string,
) (phpSetting entity.PhpSetting, err error) {
	if setting == "" {
		return phpSetting, errors.New("InvalidPhpSetting")
	}

	settingParts := strings.Split(setting, " ")
	if len(settingParts) != 2 {
		return phpSetting, errors.New("InvalidPhpSetting")
	}

	settingNameStr := settingParts[0]
	settingValueStr := settingParts[1]
	if settingNameStr == "" || settingValueStr == "" {
		return phpSetting, errors.New("InvalidPhpSetting")
	}

	settingName, err := valueObject.NewPhpSettingName(settingNameStr)
	if err != nil {
		return phpSetting, errors.New("InvalidPhpSettingName")
	}

	settingValue, err := valueObject.NewPhpSettingValue(settingValueStr)
	if err != nil {
		return phpSetting, errors.New("InvalidPhpSettingValue")
	}

	settingOptions := []valueObject.PhpSettingOption{}
	valuesToInject := []string{}

	switch settingValue.GetType() {
	case "bool":
		valuesToInject = []string{"On", "Off"}
	case "number":
		valuesToInject = []string{
			"0", "30", "60", "120", "300", "600", "900", "1800", "3600", "7200",
		}
	case "byteSize":
		lastChar := settingValue[len(settingValue)-1]
		switch lastChar {
		case 'K':
			valuesToInject = []string{"4096K", "8192K", "16384K"}
		case 'M':
			valuesToInject = []string{"16M", "32M", "64M", "128M", "256M", "512M", "1024M", "2048M"}
		case 'G':
			valuesToInject = []string{"1G", "2G", "4G"}
		}
	}

	switch settingName {
	case "error_reporting":
		valuesToInject = []string{
			"E_ALL",
			"~E_ALL",
			"E_ALL & ~E_DEPRECATED & ~E_STRICT",
			"E_ALL & ~E_DEPRECATED & ~E_STRICT & ~E_NOTICE & ~E_WARNING",
			"E_ERROR|E_CORE_ERROR|E_COMPILE_ERROR",
		}
	case "date.timezone":
		valuesToInject, err = repo.getPhpTimezones()
		if err != nil {
			log.Printf("FailedToGetPhpTimezones: %s", err.Error())
			valuesToInject = []string{}
		}
	}

	if len(valuesToInject) > 0 {
		for _, valueToInject := range valuesToInject {
			settingOption, _ := valueObject.NewPhpSettingOption(valueToInject)
			settingOptions = append(settingOptions, settingOption)
		}
	}

	settingTypeStr := "text"
	if len(settingOptions) > 0 {
		settingTypeStr = "select"
	}
	settingType, _ := valueObject.NewPhpSettingType(settingTypeStr)

	return entity.NewPhpSetting(
		settingName, settingType, settingValue, settingOptions,
	), nil
}

func (repo RuntimeQueryRepo) ReadPhpSettings(
	hostname valueObject.Fqdn,
) (phpSettings []entity.PhpSetting, err error) {
	vhostPhpConfFilePath, err := repo.GetVirtualHostPhpConfFilePath(hostname)
	if err != nil {
		return phpSettings, err
	}

	output, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "sed",
		Args: []string{
			"-n",
			"/phpIniOverride\\s*{/,/}/ { /phpIniOverride\\s*{/d; /}/d; " +
				"s/^[[:space:]]*//; s/[^[:space:]]*[[:space:]]//; p; }",
			vhostPhpConfFilePath.String(),
		},
	})
	if err != nil || output == "" {
		return phpSettings, errors.New("GetPhpSettingsFailed: " + err.Error())
	}

	for _, setting := range strings.Split(output, "\n") {
		phpSetting, err := repo.phpSettingFactory(setting)
		if err != nil {
			continue
		}

		phpSettings = append(phpSettings, phpSetting)
	}

	return phpSettings, nil
}

func (repo RuntimeQueryRepo) ReadPhpModules(
	version valueObject.PhpVersion,
) (phpModules []entity.PhpModule, err error) {
	activeModuleList, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "/usr/local/lsws/lsphp" + version.GetWithoutDots() + "/bin/php",
		Args:    []string{"-m"},
	})
	if err != nil {
		return phpModules, errors.New("GetActivePhpModulesFailed: " + err.Error())
	}

	activeModules := []string{}
	for _, moduleName := range strings.Split(activeModuleList, "\n") {
		if moduleName == "" {
			continue
		}

		moduleName = strings.Replace(moduleName, "Zend", "", -1)
		moduleName = strings.Replace(moduleName, "Loader", "", -1)

		phpModule, err := valueObject.NewPhpModuleName(moduleName)
		if err != nil {
			continue
		}

		activeModules = append(activeModules, phpModule.String())
	}

	for _, moduleName := range valueObject.ValidPhpModuleNames {
		isModuleInstalled := false
		if slices.Contains(activeModules, moduleName) {
			isModuleInstalled = true
		}

		phpModule, err := valueObject.NewPhpModuleName(moduleName)
		if err != nil {
			continue
		}

		phpModules = append(
			phpModules, entity.NewPhpModule(phpModule, isModuleInstalled),
		)
	}

	return phpModules, nil
}

func (repo RuntimeQueryRepo) ReadPhpConfigs(
	hostname valueObject.Fqdn,
) (phpConfigs entity.PhpConfigs, err error) {
	phpVersion, err := repo.ReadPhpVersion(hostname)
	if err != nil {
		return phpConfigs, err
	}

	phpSettings, err := repo.ReadPhpSettings(hostname)
	if err != nil {
		return phpConfigs, err
	}

	phpModules, err := repo.ReadPhpModules(phpVersion.Value)
	if err != nil {
		return phpConfigs, err
	}

	return entity.NewPhpConfigs(
		hostname,
		phpVersion,
		phpSettings,
		phpModules,
	), nil
}
