package runtimeInfra

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	"golang.org/x/exp/slices"
)

type RuntimeQueryRepo struct {
}

func (repo RuntimeQueryRepo) GetVirtualHostPhpConfFilePath(
	hostname valueObject.Fqdn,
) (valueObject.UnixFilePath, error) {
	var vhostPhpConfFilePath valueObject.UnixFilePath

	primaryVhostPhpConfFilePathStr := "/app/conf/php/primary.conf"
	vhostPhpConfFilePathStr := "/app/conf/php/" + hostname.String() + ".conf"

	primaryVhost, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		return vhostPhpConfFilePath, errors.New("PrimaryVhostNotFound: " + err.Error())
	}

	if hostname.String() == primaryVhost.String() {
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

func (repo RuntimeQueryRepo) GetPhpVersionsInstalled() ([]valueObject.PhpVersion, error) {
	olsConfigFile := "/usr/local/lsws/conf/httpd_config.conf"
	output, err := infraHelper.RunCmd(
		"awk",
		"/extprocessor lsphp/{print $2}",
		olsConfigFile,
	)
	if err != nil {
		return nil, errors.New("FailedToGetPhpVersions: " + err.Error())
	}

	phpVersions := []valueObject.PhpVersion{}
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

func (repo RuntimeQueryRepo) GetPhpVersion(
	hostname valueObject.Fqdn,
) (entity.PhpVersion, error) {
	var phpVersion entity.PhpVersion

	vhostPhpConfFilePath, err := repo.GetVirtualHostPhpConfFilePath(hostname)
	if err != nil {
		return phpVersion, err
	}

	currentPhpVersionStr, err := infraHelper.RunCmd(
		"awk",
		"/lsapi:lsphp/ {gsub(/[^0-9]/, \"\", $2); print $2}",
		vhostPhpConfFilePath.String(),
	)
	if err != nil {
		return phpVersion, errors.New("FailedToGetPhpVersion: " + err.Error())
	}

	currentPhpVersion, err := valueObject.NewPhpVersion(currentPhpVersionStr)
	if err != nil {
		return phpVersion, errors.New("FailedToGetPhpVersion: " + err.Error())
	}

	phpVersions, err := repo.GetPhpVersionsInstalled()
	if err != nil {
		return phpVersion, errors.New("FailedToGetPhpVersion: " + err.Error())
	}

	phpVersion = entity.NewPhpVersion(currentPhpVersion, phpVersions)
	return phpVersion, nil
}

func (repo RuntimeQueryRepo) getPhpTimezones() ([]string, error) {
	var timezones []string

	timezonesRaw, err := infraHelper.RunCmd(
		"php",
		"-r",
		"echo json_encode(DateTimeZone::listIdentifiers());",
	)
	if err != nil {
		return timezones, errors.New("FailedToGetPhpTimezones: " + err.Error())
	}

	err = json.Unmarshal([]byte(timezonesRaw), &timezones)
	if err != nil {
		return timezones, errors.New("FailedToGetPhpTimezones: " + err.Error())
	}

	return timezones, nil
}

func (repo RuntimeQueryRepo) phpSettingFactory(
	setting string,
) (entity.PhpSetting, error) {
	var phpSetting entity.PhpSetting

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
			settingOptions = append(
				settingOptions,
				valueObject.NewPhpSettingOptionPanic(valueToInject),
			)
		}
	}

	return entity.NewPhpSetting(settingName, settingValue, settingOptions), nil
}

func (repo RuntimeQueryRepo) GetPhpSettings(
	hostname valueObject.Fqdn,
) ([]entity.PhpSetting, error) {
	phpSettings := []entity.PhpSetting{}

	vhostPhpConfFilePath, err := repo.GetVirtualHostPhpConfFilePath(hostname)
	if err != nil {
		return phpSettings, err
	}

	output, err := infraHelper.RunCmd(
		"sed",
		"-n",
		"/phpIniOverride\\s*{/,/}/ { /phpIniOverride\\s*{/d; /}/d; s/^[[:space:]]*//; s/[^[:space:]]*[[:space:]]//; p; }",
		vhostPhpConfFilePath.String(),
	)
	if err != nil || output == "" {
		return phpSettings, errors.New("FailedToGetPhpSettings: " + err.Error())
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

func (repo RuntimeQueryRepo) GetPhpModules(
	version valueObject.PhpVersion,
) ([]entity.PhpModule, error) {
	phpModules := []entity.PhpModule{}

	activeModuleList, err := infraHelper.RunCmd(
		"/usr/local/lsws/lsphp"+version.GetWithoutDots()+"/bin/php",
		"-m",
	)
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
			phpModules,
			entity.NewPhpModule(phpModule, isModuleInstalled),
		)
	}

	return phpModules, nil
}

func (repo RuntimeQueryRepo) GetPhpConfigs(
	hostname valueObject.Fqdn,
) (entity.PhpConfigs, error) {
	var phpConfigs entity.PhpConfigs

	phpVersion, err := repo.GetPhpVersion(hostname)
	if err != nil {
		return phpConfigs, err
	}

	phpSettings, err := repo.GetPhpSettings(hostname)
	if err != nil {
		return phpConfigs, err
	}

	phpModules, err := repo.GetPhpModules(phpVersion.Value)
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
