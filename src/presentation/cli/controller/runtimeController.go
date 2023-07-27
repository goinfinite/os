package cliController

import (
	"os"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/domain/valueObject"
	"github.com/speedianet/sam/src/infra"
	cliHelper "github.com/speedianet/sam/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetPhpConfigsController() *cobra.Command {
	var hostnameStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetPhpConfigs",
		Run: func(cmd *cobra.Command, args []string) {
			if hostnameStr == "" {
				hostnameStr = os.Getenv("VIRTUAL_HOST")
			}
			hostname := valueObject.NewFqdnPanic(hostnameStr)

			runtimeQueryRepo := infra.RuntimeQueryRepo{}
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
			if hostnameStr == "" {
				hostnameStr = os.Getenv("VIRTUAL_HOST")
			}
			hostname := valueObject.NewFqdnPanic(hostnameStr)

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

			if len(phpModules) == 0 && len(phpSettings) == 0 {
				cliHelper.ResponseWrapper(false, "MissingModuleOrSetting")
			}

			updatePhpConfigsDto := dto.NewUpdatePhpConfigs(
				hostname,
				phpVersion,
				phpModules,
				phpSettings,
			)

			runtimeQueryRepo := infra.RuntimeQueryRepo{}
			runtimeCmdRepo := infra.RuntimeCmdRepo{}
			wsQueryRepo := infra.WsQueryRepo{}

			err := useCase.UpdatePhpConfigs(
				runtimeQueryRepo,
				runtimeCmdRepo,
				wsQueryRepo,
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
		Use:   "updateSetting",
		Short: "UpdatePhpSetting",
		Run: func(cmd *cobra.Command, args []string) {
			if hostnameStr == "" {
				hostnameStr = os.Getenv("VIRTUAL_HOST")
			}
			hostname := valueObject.NewFqdnPanic(hostnameStr)

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

			runtimeQueryRepo := infra.RuntimeQueryRepo{}
			runtimeCmdRepo := infra.RuntimeCmdRepo{}
			wsQueryRepo := infra.WsQueryRepo{}

			err := useCase.UpdatePhpConfigs(
				runtimeQueryRepo,
				runtimeCmdRepo,
				wsQueryRepo,
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
		Use:   "updateModule",
		Short: "UpdatePhpModule",
		Run: func(cmd *cobra.Command, args []string) {
			if hostnameStr == "" {
				hostnameStr = os.Getenv("VIRTUAL_HOST")
			}
			hostname := valueObject.NewFqdnPanic(hostnameStr)

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

			runtimeQueryRepo := infra.RuntimeQueryRepo{}
			runtimeCmdRepo := infra.RuntimeCmdRepo{}
			wsQueryRepo := infra.WsQueryRepo{}

			err := useCase.UpdatePhpConfigs(
				runtimeQueryRepo,
				runtimeCmdRepo,
				wsQueryRepo,
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
