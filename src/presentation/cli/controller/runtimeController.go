package cliController

import (
	"errors"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	runtimeInfra "github.com/speedianet/os/src/infra/runtime"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	sharedHelper "github.com/speedianet/os/src/presentation/shared/helper"
	"github.com/spf13/cobra"
)

func getHostname(hostnameStr string) (valueObject.Fqdn, error) {
	primaryVhost, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		return "", errors.New("PrimaryVirtualHostNotFound")
	}

	hostname := primaryVhost
	if hostnameStr != "" {
		hostname = valueObject.NewFqdnPanic(hostnameStr)
	}

	return hostname, nil
}

func GetPhpConfigsController() *cobra.Command {
	var hostnameStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetPhpConfigs",
		Run: func(cmd *cobra.Command, args []string) {
			svcName := valueObject.NewServiceNamePanic("php")
			sharedHelper.StopIfServiceUnavailable(svcName.String())

			hostname, err := getHostname(hostnameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			runtimeQueryRepo := runtimeInfra.RuntimeQueryRepo{}
			phpConfigs, err := useCase.GetPhpConfigs(runtimeQueryRepo, hostname)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, phpConfigs)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "H", "", "Hostname")
	return cmd
}

func UpdatePhpConfigController() *cobra.Command {
	var hostnameStr string
	var phpVersionStr string
	var moduleNameStr string
	moduleStatusBool := true
	var settingNameStr string
	var settingValueStr string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdatePhpConfigs",
		Run: func(cmd *cobra.Command, args []string) {
			svcName := valueObject.NewServiceNamePanic("php")
			sharedHelper.StopIfServiceUnavailable(svcName.String())

			hostname, err := getHostname(hostnameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			phpVersion := valueObject.NewPhpVersionPanic(phpVersionStr)

			phpModules := []entity.PhpModule{}
			if moduleNameStr != "" {
				moduleName := valueObject.NewPhpModuleNamePanic(moduleNameStr)
				phpModules = append(
					phpModules,
					entity.NewPhpModule(moduleName, moduleStatusBool),
				)
			}

			phpSettings := []entity.PhpSetting{}
			if settingNameStr != "" {
				settingName := valueObject.NewPhpSettingNamePanic(settingNameStr)
				settingValue := valueObject.NewPhpSettingValuePanic(settingValueStr)
				phpSettings = append(
					phpSettings,
					entity.NewPhpSetting(settingName, settingValue, nil),
				)
			}

			updatePhpConfigsDto := dto.NewUpdatePhpConfigs(
				hostname,
				phpVersion,
				phpModules,
				phpSettings,
			)

			runtimeQueryRepo := runtimeInfra.RuntimeQueryRepo{}
			runtimeCmdRepo := runtimeInfra.NewRuntimeCmdRepo()
			vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}

			err = useCase.UpdatePhpConfigs(
				runtimeQueryRepo,
				runtimeCmdRepo,
				vhostQueryRepo,
				updatePhpConfigsDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "PhpConfigsUpdated")
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "H", "", "Hostname")
	cmd.Flags().StringVarP(&phpVersionStr, "version", "v", "", "PhpVersion")
	cmd.MarkFlagRequired("version")
	cmd.Flags().StringVarP(&moduleNameStr, "module", "m", "", "PhpModuleName")
	cmd.Flags().BoolVarP(&moduleStatusBool, "status", "s", true, "PhpModuleStatus")
	cmd.Flags().StringVarP(&settingNameStr, "setting", "S", "", "PhpSettingName")
	cmd.Flags().StringVarP(&settingValueStr, "value", "V", "", "PhpSettingValue")
	return cmd
}

func UpdatePhpSettingController() *cobra.Command {
	var hostnameStr string
	var phpVersionStr string
	var settingNameStr string
	var settingValueStr string

	cmd := &cobra.Command{
		Use:   "update-setting",
		Short: "UpdatePhpSetting",
		Run: func(cmd *cobra.Command, args []string) {
			svcName := valueObject.NewServiceNamePanic("php")
			sharedHelper.StopIfServiceUnavailable(svcName.String())

			hostname, err := getHostname(hostnameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			phpVersion := valueObject.NewPhpVersionPanic(phpVersionStr)

			phpSettings := []entity.PhpSetting{}
			settingName := valueObject.NewPhpSettingNamePanic(settingNameStr)
			settingValue := valueObject.NewPhpSettingValuePanic(settingValueStr)
			phpSettings = append(
				phpSettings,
				entity.NewPhpSetting(settingName, settingValue, nil),
			)

			phpModules := []entity.PhpModule{}

			updatePhpConfigsDto := dto.NewUpdatePhpConfigs(
				hostname,
				phpVersion,
				phpModules,
				phpSettings,
			)

			runtimeQueryRepo := runtimeInfra.RuntimeQueryRepo{}
			runtimeCmdRepo := runtimeInfra.NewRuntimeCmdRepo()
			vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}

			err = useCase.UpdatePhpConfigs(
				runtimeQueryRepo,
				runtimeCmdRepo,
				vhostQueryRepo,
				updatePhpConfigsDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "PhpSettingUpdated")
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "H", "", "Hostname")
	cmd.Flags().StringVarP(&phpVersionStr, "version", "v", "", "PhpVersion")
	cmd.MarkFlagRequired("version")
	cmd.Flags().StringVarP(&settingNameStr, "name", "n", "", "PhpSettingName")
	cmd.Flags().StringVarP(&settingValueStr, "value", "V", "", "PhpSettingValue")
	return cmd
}

func UpdatePhpModuleController() *cobra.Command {
	var hostnameStr string
	var phpVersionStr string
	var moduleNameStr string
	moduleStatusBool := true

	cmd := &cobra.Command{
		Use:   "update-module",
		Short: "UpdatePhpModule",
		Run: func(cmd *cobra.Command, args []string) {
			svcName := valueObject.NewServiceNamePanic("php")
			sharedHelper.StopIfServiceUnavailable(svcName.String())

			hostname, err := getHostname(hostnameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			phpVersion := valueObject.NewPhpVersionPanic(phpVersionStr)

			phpModules := []entity.PhpModule{}
			moduleName := valueObject.NewPhpModuleNamePanic(moduleNameStr)
			phpModules = append(
				phpModules,
				entity.NewPhpModule(moduleName, moduleStatusBool),
			)

			phpSettings := []entity.PhpSetting{}

			updatePhpConfigsDto := dto.NewUpdatePhpConfigs(
				hostname,
				phpVersion,
				phpModules,
				phpSettings,
			)

			runtimeQueryRepo := runtimeInfra.RuntimeQueryRepo{}
			runtimeCmdRepo := runtimeInfra.NewRuntimeCmdRepo()
			vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}

			err = useCase.UpdatePhpConfigs(
				runtimeQueryRepo,
				runtimeCmdRepo,
				vhostQueryRepo,
				updatePhpConfigsDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "PhpModuleUpdated")
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "H", "", "Hostname")
	cmd.Flags().StringVarP(&phpVersionStr, "version", "v", "", "PhpVersion")
	cmd.MarkFlagRequired("version")
	cmd.Flags().StringVarP(&moduleNameStr, "name", "n", "", "PhpModuleName")
	cmd.Flags().BoolVarP(&moduleStatusBool, "status", "s", true, "PhpModuleStatus")
	return cmd
}
