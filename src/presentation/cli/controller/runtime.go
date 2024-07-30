package cliController

import (
	"errors"
	"strconv"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/speedianet/os/src/presentation/service"
	sharedHelper "github.com/speedianet/os/src/presentation/shared/helper"
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

func (controller *RuntimeController) parsePhpModules(
	rawPhpModules []string,
) []entity.PhpModule {
	modules := []entity.PhpModule{}
	if len(rawPhpModules) == 0 {
		return modules
	}

	for _, rawModule := range rawPhpModules {
		rawModuleParts := strings.Split(rawModule, ":")
		rawModulePartsLength := len(rawModuleParts)
		if rawModulePartsLength == 0 {
			continue
		}

		moduleName, err := valueObject.NewPhpModuleName(rawModuleParts[0])
		if err != nil {
			continue
		}

		moduleStatus := true
		if rawModulePartsLength > 1 {
			moduleStatus, err = sharedHelper.ParseBoolParam(rawModuleParts[1])
			if err != nil {
				moduleStatus = false
			}
		}

		modules = append(modules, entity.NewPhpModule(moduleName, moduleStatus))
	}

	return modules
}

func (controller *RuntimeController) parsePhpSettings(
	rawPhpSettings []string,
) []entity.PhpSetting {
	settings := []entity.PhpSetting{}
	if len(rawPhpSettings) == 0 {
		return settings
	}

	for _, rawSetting := range rawPhpSettings {
		rawSettingParts := strings.Split(rawSetting, ":")
		if len(rawSettingParts) != 2 {
			continue
		}

		settingName, err := valueObject.NewPhpSettingName(rawSettingParts[0])
		if err != nil {
			continue
		}

		settingValue, err := valueObject.NewPhpSettingValue(rawSettingParts[1])
		if err != nil {
			continue
		}

		settings = append(settings, entity.NewPhpSetting(settingName, settingValue, nil))
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

			requestBody["modules"] = controller.parsePhpModules(modulesSlice)
			requestBody["settings"] = controller.parsePhpSettings(settingsSlice)

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
			requestBody["modules"] = controller.parsePhpModules(
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
			requestBody["settings"] = controller.parsePhpSettings(
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
