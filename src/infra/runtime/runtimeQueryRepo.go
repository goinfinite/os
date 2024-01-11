package runtimeInfra

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	"github.com/speedianet/os/src/infra"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	"golang.org/x/exp/slices"
)

type RuntimeQueryRepo struct {
}

func (r RuntimeQueryRepo) GetPhpVersionsInstalled() ([]valueObject.PhpVersion, error) {
	olsConfigFile := "/usr/local/lsws/conf/httpd_config.conf"
	output, err := infraHelper.RunCmd(
		"awk",
		"/extprocessor lsphp/{print $2}",
		olsConfigFile,
	)
	if err != nil {
		log.Printf("FailedToGetPhpVersions: %v", err)
		return nil, errors.New("FailedToGetPhpVersions")
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

func (r RuntimeQueryRepo) GetPhpVersion(
	hostname valueObject.Fqdn,
) (entity.PhpVersion, error) {
	vhconfFile := infra.WsQueryRepo{}.GetVirtualHostConfFilePath(hostname)
	currentPhpVersionStr, err := infraHelper.RunCmd(
		"awk",
		"/lsapi:lsphp/ {gsub(/[^0-9]/, \"\", $2); print $2}",
		vhconfFile,
	)
	if err != nil {
		log.Printf("FailedToGetPhpVersion: %v", err)
		return entity.PhpVersion{}, errors.New("FailedToGetPhpVersion")
	}

	currentPhpVersion, err := valueObject.NewPhpVersion(currentPhpVersionStr)
	if err != nil {
		return entity.PhpVersion{}, errors.New("FailedToGetPhpVersion")
	}

	phpVersions, err := r.GetPhpVersionsInstalled()
	if err != nil {
		return entity.PhpVersion{}, errors.New("FailedToGetPhpVersion")
	}

	return entity.NewPhpVersion(currentPhpVersion, phpVersions), nil
}

func (r RuntimeQueryRepo) getPhpTimezones() ([]string, error) {
	timezonesRaw, err := infraHelper.RunCmd(
		"php",
		"-r",
		"echo json_encode(DateTimeZone::listIdentifiers());",
	)
	if err != nil {
		log.Printf("FailedToGetPhpTimezones: %v", err)
		return nil, errors.New("FailedToGetPhpTimezones")
	}

	var timezones []string
	err = json.Unmarshal([]byte(timezonesRaw), &timezones)
	if err != nil {
		log.Printf("FailedToGetPhpTimezones: %v", err)
		return nil, errors.New("FailedToGetPhpTimezones")
	}

	return timezones, nil
}

func (r RuntimeQueryRepo) phpSettingFactory(
	setting string,
) (entity.PhpSetting, error) {
	if setting == "" {
		return entity.PhpSetting{}, errors.New("InvalidPhpSetting")
	}

	settingParts := strings.Split(setting, " ")
	if len(settingParts) != 2 {
		return entity.PhpSetting{}, errors.New("InvalidPhpSetting")
	}

	settingNameStr := settingParts[0]
	settingValueStr := settingParts[1]
	if settingNameStr == "" || settingValueStr == "" {
		return entity.PhpSetting{}, errors.New("InvalidPhpSetting")
	}

	settingName, err := valueObject.NewPhpSettingName(settingNameStr)
	if err != nil {
		return entity.PhpSetting{}, errors.New("InvalidPhpSettingName")
	}

	settingValue, err := valueObject.NewPhpSettingValue(settingValueStr)
	if err != nil {
		return entity.PhpSetting{}, errors.New("InvalidPhpSettingValue")
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
		valuesToInject, err = r.getPhpTimezones()
		if err != nil {
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

func (r RuntimeQueryRepo) GetPhpSettings(
	hostname valueObject.Fqdn,
) ([]entity.PhpSetting, error) {
	vhconfFile := infra.WsQueryRepo{}.GetVirtualHostConfFilePath(hostname)
	output, err := infraHelper.RunCmd(
		"sed",
		"-n",
		"/phpIniOverride\\s*{/,/}/ { /phpIniOverride\\s*{/d; /}/d; s/^[[:space:]]*//; s/[^[:space:]]*[[:space:]]//; p; }",
		vhconfFile,
	)
	if err != nil || output == "" {
		log.Printf("FailedToGetPhpSettings: %v", err)
		return nil, errors.New("FailedToGetPhpSettings")
	}

	phpSettings := []entity.PhpSetting{}
	for _, setting := range strings.Split(output, "\n") {
		phpSetting, err := r.phpSettingFactory(setting)
		if err != nil {
			continue
		}

		phpSettings = append(phpSettings, phpSetting)
	}

	return phpSettings, nil
}

func (r RuntimeQueryRepo) GetPhpModules(
	version valueObject.PhpVersion,
) ([]entity.PhpModule, error) {
	activeModuleList, err := infraHelper.RunCmd(
		"/usr/local/lsws/lsphp"+version.GetWithoutDots()+"/bin/php",
		"-m",
	)
	if err != nil {
		log.Printf("GetActivePhpModulesFailed: %v", err)
		return nil, errors.New("GetActivePhpModulesFailed")
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

	phpModules := []entity.PhpModule{}
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

func (r RuntimeQueryRepo) GetPhpConfigs(
	hostname valueObject.Fqdn,
) (entity.PhpConfigs, error) {
	phpVersion, err := r.GetPhpVersion(hostname)
	if err != nil {
		return entity.PhpConfigs{}, err
	}

	phpSettings, err := r.GetPhpSettings(hostname)
	if err != nil {
		return entity.PhpConfigs{}, err
	}

	phpModules, err := r.GetPhpModules(phpVersion.Value)
	if err != nil {
		return entity.PhpConfigs{}, err
	}

	return entity.NewPhpConfigs(
		hostname,
		phpVersion,
		phpSettings,
		phpModules,
	), nil
}
