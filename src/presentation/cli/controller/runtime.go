package cliController

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/speedianet/os/src/presentation/service"
	"github.com/spf13/cobra"
)

type RuntimeController struct {
	runtimeService *service.RuntimeService
}

func NewRuntimeController(
	persistentDbService *internalDbInfra.PersistentDatabaseService,
) *RuntimeController {
	return &RuntimeController{
		runtimeService: service.NewRuntimeService(persistentDbService),
	}
}

func getHostname(hostnameStr string) (hostname valueObject.Fqdn, err error) {
	primaryVhost, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		return hostname, errors.New("PrimaryVirtualHostNotFound")
	}

	hostname = primaryVhost
	if hostnameStr != "" {
		return valueObject.NewFqdn(hostnameStr)
	}

	return hostname, nil
}

func (controller *RuntimeController) ReadPhpConfigs() *cobra.Command {
	var hostnameStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetPhpConfigs",
		Run: func(cmd *cobra.Command, args []string) {
			hostname, err := getHostname(hostnameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			requestBody := map[string]interface{}{
				"hostname": hostname.String(),
			}

			cliHelper.ServiceResponseWrapper(
				controller.runtimeService.ReadPhpConfigs(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	return cmd
}

func parsePhpModules(rawPhpModules []string) []entity.PhpModule {
	modules := []entity.PhpModule{}
	for _, rawModule := range rawPhpModules {
		rawModuleParts := strings.Split(rawModule, ":")
		rawModulePartsLength := len(rawModuleParts)
		if rawModulePartsLength == 0 {
			slog.Debug("PhpModuleEmpty", slog.String("module", rawModule))
			continue
		}

		moduleName, err := valueObject.NewPhpModuleName(rawModuleParts[0])
		if err != nil {
			slog.Debug("InvalidPhpModuleName", slog.Any("name", rawModuleParts[0]))
			continue
		}

		moduleStatus := true
		if rawModulePartsLength > 1 {
			var err error
			moduleStatus, err = voHelper.InterfaceToBool(rawModuleParts[1])
			if err != nil {
				moduleStatus = false
			}
		}

		modules = append(modules, entity.NewPhpModule(moduleName, moduleStatus))
	}

	return modules
}

func parsePhpSettings(rawPhpSettings []string) []entity.PhpSetting {
	settings := []entity.PhpSetting{}
	for _, rawSetting := range rawPhpSettings {
		rawSettingParts := strings.Split(rawSetting, ":")
		if len(rawSettingParts) == 0 {
			slog.Debug("PhpSettingEmpty", slog.String("setting", rawSetting))
			continue
		}

		settingName, err := valueObject.NewPhpSettingName(rawSettingParts[0])
		if err != nil {
			slog.Debug(
				"InvalidPhpSettingName", slog.Any("name", rawSettingParts[0]),
			)
			continue
		}

		settingValue, err := valueObject.NewPhpSettingValue(rawSettingParts[1])
		if err != nil {
			slog.Debug(
				"InvalidPhpSettingValue", slog.Any("value", rawSettingParts[1]),
			)
			continue
		}

		settings = append(
			settings,
			entity.NewPhpSetting(
				settingName, settingValue, []valueObject.PhpSettingOption{},
			),
		)
	}

	return settings
}

func (controller *RuntimeController) UpdatePhpConfig() *cobra.Command {
	var hostnameStr, phpVersionStr string
	var modulesSlice, settingsSlice []string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdatePhpConfigs",
		Run: func(cmd *cobra.Command, args []string) {
			hostname, err := getHostname(hostnameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			requestBody := map[string]interface{}{
				"hostname": hostname.String(),
				"version":  phpVersionStr,
			}

			if len(modulesSlice) > 0 {
				requestBody["modules"] = parsePhpModules(modulesSlice)
			}

			if len(settingsSlice) > 0 {
				requestBody["settings"] = parsePhpSettings(settingsSlice)
			}

			cliHelper.ServiceResponseWrapper(
				controller.runtimeService.UpdatePhpConfigs(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.Flags().StringVarP(&phpVersionStr, "version", "v", "", "PhpVersion")
	cmd.MarkFlagRequired("version")
	cmd.Flags().StringSliceVarP(
		&modulesSlice, "module", "m", []string{}, "(phpModuleName:phpModuleStatus)",
	)
	cmd.Flags().StringSliceVarP(
		&settingsSlice, "setting", "s", []string{}, "(phpSettingName:phpSettingValue)",
	)
	return cmd
}

func (controller *RuntimeController) UpdatePhpModule() *cobra.Command {
	var hostnameStr, phpVersionStr, moduleNameStr string
	moduleStatusBool := true

	cmd := &cobra.Command{
		Use:   "update-module",
		Short: "UpdatePhpModule",
		Run: func(cmd *cobra.Command, args []string) {
			hostname, err := getHostname(hostnameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}
			requestBody := map[string]interface{}{
				"hostname": hostname.String(),
				"version":  phpVersionStr,
			}

			moduleStatusStr := strconv.FormatBool(moduleStatusBool)
			rawPhpModuleParam := moduleNameStr + ":" + moduleStatusStr
			requestBody["modules"] = parsePhpModules(
				[]string{rawPhpModuleParam},
			)

			cliHelper.ServiceResponseWrapper(
				controller.runtimeService.UpdatePhpConfigs(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.Flags().StringVarP(&phpVersionStr, "version", "v", "", "PhpVersion")
	cmd.MarkFlagRequired("version")
	cmd.Flags().StringVarP(&moduleNameStr, "name", "N", "", "PhpModuleName")
	cmd.MarkFlagRequired("name")
	cmd.Flags().BoolVarP(&moduleStatusBool, "status", "V", true, "PhpModuleStatus")
	cmd.MarkFlagRequired("status")
	return cmd
}

func (controller *RuntimeController) UpdatePhpSetting() *cobra.Command {
	var hostnameStr, phpVersionStr, settingNameStr, settingValueStr string

	cmd := &cobra.Command{
		Use:   "update-setting",
		Short: "UpdatePhpSetting",
		Run: func(cmd *cobra.Command, args []string) {
			hostname, err := getHostname(hostnameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}
			requestBody := map[string]interface{}{
				"hostname": hostname.String(),
				"version":  phpVersionStr,
			}

			rawPhpSettingParam := settingNameStr + ":" + settingValueStr
			requestBody["settings"] = parsePhpSettings(
				[]string{rawPhpSettingParam},
			)

			cliHelper.ServiceResponseWrapper(
				controller.runtimeService.UpdatePhpConfigs(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.Flags().StringVarP(&phpVersionStr, "version", "v", "", "PhpVersion")
	cmd.MarkFlagRequired("version")
	cmd.Flags().StringVarP(&settingNameStr, "name", "N", "", "PhpSettingName")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&settingValueStr, "value", "V", "", "PhpSettingValue")
	cmd.MarkFlagRequired("value")
	return cmd
}
