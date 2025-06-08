package cliController

import (
	"errors"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	cliHelper "github.com/goinfinite/os/src/presentation/cli/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
	"github.com/spf13/cobra"
)

type RuntimeController struct {
	runtimeLiaison *liaison.RuntimeLiaison
}

func NewRuntimeController(
	persistentDbService *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *RuntimeController {
	return &RuntimeController{
		runtimeLiaison: liaison.NewRuntimeLiaison(persistentDbService, trailDbSvc),
	}
}

func getHostname(hostnameStr string) (hostname valueObject.Fqdn, err error) {
	primaryVhost, err := infraHelper.ReadPrimaryVirtualHostHostname()
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

			cliHelper.LiaisonResponseWrapper(
				controller.runtimeLiaison.ReadPhpConfigs(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	return cmd
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
				modules := []entity.PhpModule{}
				for _, rawModule := range modulesSlice {
					module, err := entity.NewPhpModuleFromString(rawModule)
					if err != nil {
						continue
					}
					modules = append(modules, module)
				}
				requestBody["modules"] = modules
			}

			if len(settingsSlice) > 0 {
				settings := []entity.PhpSetting{}
				for _, rawSetting := range settingsSlice {
					setting, err := entity.NewPhpSettingFromString(rawSetting)
					if err != nil {
						continue
					}
					settings = append(settings, setting)
				}
				requestBody["settings"] = settings
			}

			cliHelper.LiaisonResponseWrapper(
				controller.runtimeLiaison.UpdatePhpConfigs(requestBody),
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
	var hostnameStr, phpVersionStr, moduleNameStr, moduleStatusStr string

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

			rawPhpModuleParam := moduleNameStr + ":" + moduleStatusStr
			module, err := entity.NewPhpModuleFromString(rawPhpModuleParam)
			if err != nil {
				cliHelper.ResponseWrapper(false, err)
			}
			requestBody["modules"] = []entity.PhpModule{module}

			cliHelper.LiaisonResponseWrapper(
				controller.runtimeLiaison.UpdatePhpConfigs(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.Flags().StringVarP(&phpVersionStr, "version", "v", "", "PhpVersion")
	cmd.MarkFlagRequired("version")
	cmd.Flags().StringVarP(&moduleNameStr, "name", "N", "", "PhpModuleName")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&moduleStatusStr, "status", "V", "true", "PhpModuleStatus")
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
			setting, err := entity.NewPhpSettingFromString(rawPhpSettingParam)
			if err != nil {
				cliHelper.ResponseWrapper(false, err)
			}
			requestBody["settings"] = []entity.PhpSetting{setting}

			cliHelper.LiaisonResponseWrapper(
				controller.runtimeLiaison.UpdatePhpConfigs(requestBody),
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

func (controller *RuntimeController) RunPhpCommand() *cobra.Command {
	var hostnameStr, commandStr string

	cmd := &cobra.Command{
		Use:   "run",
		Short: "RunPhpCommand",
		Run: func(cmd *cobra.Command, args []string) {
			hostname, err := getHostname(hostnameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}
			requestBody := map[string]interface{}{
				"hostname": hostname.String(),
				"command":  commandStr,
			}

			cliHelper.LiaisonResponseWrapper(
				controller.runtimeLiaison.RunPhpCommand(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.MarkFlagRequired("hostname")
	cmd.Flags().StringVarP(&commandStr, "command", "c", "", "Command")
	cmd.MarkFlagRequired("command")
	return cmd
}
