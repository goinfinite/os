package cliController

import (
	"errors"

	"github.com/goinfinite/os/src/domain/entity"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	cliHelper "github.com/goinfinite/os/src/presentation/cli/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
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

func getHostname(hostnameStr string) (hostname tkValueObject.Fqdn, err error) {
	primaryVhost, err := vhostInfra.NewVirtualHostHelpers().
		ReadPrimaryVirtualHostHostname()
	if err != nil {
		return hostname, errors.New("PrimaryVirtualHostNotFound")
	}

	hostname = primaryVhost
	if hostnameStr != "" {
		return tkValueObject.NewFqdn(hostnameStr)
	}

	return hostname, nil
}

func (controller *RuntimeController) ReadPhpConfigs() *cobra.Command {
	var hostnameStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetPhpConfigs",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			hostname, err := getHostname(hostnameStr)
			if err != nil {
				tkPresentation.SimpleCliResponseRenderer(false, err.Error())
				return
			}

			requestBody := map[string]interface{}{
				"hostname": hostname.String(),
			}

			tkPresentation.LiaisonCliResponseRenderer(
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
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			hostname, err := getHostname(hostnameStr)
			if err != nil {
				tkPresentation.SimpleCliResponseRenderer(false, err.Error())
				return
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
						tkPresentation.SimpleCliResponseRenderer(false, err.Error())
						return
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
						tkPresentation.SimpleCliResponseRenderer(false, err.Error())
						return
					}
					settings = append(settings, setting)
				}
				requestBody["settings"] = settings
			}

			tkPresentation.LiaisonCliResponseRenderer(
				controller.runtimeLiaison.UpdatePhpConfigs(requestBody, false),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.Flags().StringVarP(&phpVersionStr, "version", "v", "", "PhpVersion")
	cliHelper.MarkRequiredFlag(cmd, "version")
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
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			hostname, err := getHostname(hostnameStr)
			if err != nil {
				tkPresentation.SimpleCliResponseRenderer(false, err.Error())
				return
			}
			requestBody := map[string]interface{}{
				"hostname": hostname.String(),
				"version":  phpVersionStr,
			}

			rawPhpModuleParam := moduleNameStr + ":" + moduleStatusStr
			module, err := entity.NewPhpModuleFromString(rawPhpModuleParam)
			if err != nil {
				tkPresentation.SimpleCliResponseRenderer(false, err.Error())
				return
			}
			requestBody["modules"] = []entity.PhpModule{module}

			tkPresentation.LiaisonCliResponseRenderer(
				controller.runtimeLiaison.UpdatePhpConfigs(requestBody, false),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.Flags().StringVarP(&phpVersionStr, "version", "v", "", "PhpVersion")
	cliHelper.MarkRequiredFlag(cmd, "version")
	cmd.Flags().StringVarP(&moduleNameStr, "name", "N", "", "PhpModuleName")
	cliHelper.MarkRequiredFlag(cmd, "name")
	cmd.Flags().StringVarP(&moduleStatusStr, "status", "V", "true", "PhpModuleStatus")
	cliHelper.MarkRequiredFlag(cmd, "status")
	return cmd
}

func (controller *RuntimeController) UpdatePhpModules() *cobra.Command {
	var hostnameStr, phpVersionStr string
	var operatorAccountIdStr, operatorIpAddressStr string
	var modulesSlice []string

	cmd := &cobra.Command{
		Use:   "update-modules",
		Short: "UpdatePhpModules",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			hostname, err := getHostname(hostnameStr)
			if err != nil {
				tkPresentation.SimpleCliResponseRenderer(false, err.Error())
				return
			}

			modules := []entity.PhpModule{}
			for _, rawModule := range modulesSlice {
				module, err := entity.NewPhpModuleFromString(rawModule)
				if err != nil {
					tkPresentation.SimpleCliResponseRenderer(false, err.Error())
					return
				}
				modules = append(modules, module)
			}

			requestBody := map[string]interface{}{
				"hostname": hostname.String(),
				"version":  phpVersionStr,
				"modules":  modules,
			}
			if operatorAccountIdStr != "" {
				requestBody["operatorAccountId"] = operatorAccountIdStr
			}
			if operatorIpAddressStr != "" {
				requestBody["operatorIpAddress"] = operatorIpAddressStr
			}

			tkPresentation.LiaisonCliResponseRenderer(
				controller.runtimeLiaison.UpdatePhpModules(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.Flags().StringVarP(&phpVersionStr, "version", "v", "", "PhpVersion")
	cliHelper.MarkRequiredFlag(cmd, "version")
	cmd.Flags().StringSliceVarP(
		&modulesSlice, "module", "m", []string{}, "(phpModuleName:phpModuleStatus)",
	)
	cliHelper.MarkRequiredFlag(cmd, "module")
	cmd.Flags().StringVar(
		&operatorAccountIdStr, "operator-account-id", "", "OperatorAccountId",
	)
	cmd.Flags().StringVar(
		&operatorIpAddressStr, "operator-ip-address", "", "OperatorIpAddress",
	)
	return cmd
}

func (controller *RuntimeController) UpdatePhpSetting() *cobra.Command {
	var hostnameStr, phpVersionStr, settingNameStr, settingValueStr string

	cmd := &cobra.Command{
		Use:   "update-setting",
		Short: "UpdatePhpSetting",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			hostname, err := getHostname(hostnameStr)
			if err != nil {
				tkPresentation.SimpleCliResponseRenderer(false, err.Error())
				return
			}
			requestBody := map[string]interface{}{
				"hostname": hostname.String(),
				"version":  phpVersionStr,
			}

			rawPhpSettingParam := settingNameStr + ":" + settingValueStr
			setting, err := entity.NewPhpSettingFromString(rawPhpSettingParam)
			if err != nil {
				tkPresentation.SimpleCliResponseRenderer(false, err.Error())
				return
			}
			requestBody["settings"] = []entity.PhpSetting{setting}

			tkPresentation.LiaisonCliResponseRenderer(
				controller.runtimeLiaison.UpdatePhpConfigs(requestBody, false),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.Flags().StringVarP(&phpVersionStr, "version", "v", "", "PhpVersion")
	cliHelper.MarkRequiredFlag(cmd, "version")
	cmd.Flags().StringVarP(&settingNameStr, "name", "N", "", "PhpSettingName")
	cliHelper.MarkRequiredFlag(cmd, "name")
	cmd.Flags().StringVarP(&settingValueStr, "value", "V", "", "PhpSettingValue")
	cliHelper.MarkRequiredFlag(cmd, "value")
	return cmd
}

func (controller *RuntimeController) RunPhpCommand() *cobra.Command {
	var hostnameStr, commandStr string

	cmd := &cobra.Command{
		Use:   "run",
		Short: "RunPhpCommand",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			hostname, err := getHostname(hostnameStr)
			if err != nil {
				tkPresentation.SimpleCliResponseRenderer(false, err.Error())
				return
			}
			requestBody := map[string]interface{}{
				"hostname": hostname.String(),
				"command":  commandStr,
			}

			tkPresentation.LiaisonCliResponseRenderer(
				controller.runtimeLiaison.RunPhpCommand(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cliHelper.MarkRequiredFlag(cmd, "hostname")
	cmd.Flags().StringVarP(&commandStr, "command", "c", "", "Command")
	cliHelper.MarkRequiredFlag(cmd, "command")
	return cmd
}
