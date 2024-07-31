package cliController

import (
	"errors"
	"strings"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	runtimeInfra "github.com/speedianet/os/src/infra/runtime"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	sharedHelper "github.com/speedianet/os/src/presentation/shared/helper"
	"github.com/spf13/cobra"
)

type RuntimeController struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewRuntimeController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *RuntimeController {
	return &RuntimeController{
		persistentDbSvc: persistentDbSvc,
	}
}

func getHostname(hostnameStr string) (valueObject.Fqdn, error) {
	primaryVhost, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		return "", errors.New("PrimaryVirtualHostNotFound")
	}

	hostname := primaryVhost
	if hostnameStr != "" {
		hostname, err = valueObject.NewFqdn(hostnameStr)
		if err != nil {
			cliHelper.ResponseWrapper(false, err.Error())
		}
	}

	return hostname, nil
}

func (controller *RuntimeController) ReadPhpConfigs() *cobra.Command {
	var hostnameStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetPhpConfigs",
		Run: func(cmd *cobra.Command, args []string) {
			serviceName, _ := valueObject.NewServiceName("php-webserver")
			sharedHelper.StopIfServiceUnavailable(controller.persistentDbSvc, serviceName)

			hostname, err := getHostname(hostnameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			runtimeQueryRepo := runtimeInfra.RuntimeQueryRepo{}
			phpConfigs, err := useCase.ReadPhpConfigs(runtimeQueryRepo, hostname)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, phpConfigs)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	return cmd
}

func (controller *RuntimeController) UpdatePhpConfig() *cobra.Command {
	var hostnameStr string
	var phpVersionStr string
	var modulesSlice []string
	var settingsSlice []string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdatePhpConfigs",
		Run: func(cmd *cobra.Command, args []string) {
			serviceName, _ := valueObject.NewServiceName("php-webserver")
			sharedHelper.StopIfServiceUnavailable(controller.persistentDbSvc, serviceName)

			phpVersion, err := valueObject.NewPhpVersion(phpVersionStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			hostname, err := getHostname(hostnameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			phpModules := []entity.PhpModule{}
			for _, rawModule := range modulesSlice {
				moduleParts := strings.Split(rawModule, ":")
				modulePartsLength := len(moduleParts)
				if modulePartsLength == 0 {
					continue
				}

				moduleName, err := valueObject.NewPhpModuleName(moduleParts[0])
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				moduleStatus := true
				if modulePartsLength > 1 {
					moduleStatus, err = sharedHelper.ParseBoolParam(moduleParts[1])
					if err != nil {
						moduleStatus = false
					}
				}

				phpModules = append(
					phpModules,
					entity.NewPhpModule(moduleName, moduleStatus),
				)
			}

			phpSettings := []entity.PhpSetting{}
			for _, rawSetting := range settingsSlice {
				settingParts := strings.Split(rawSetting, ":")
				if len(settingParts) != 2 {
					continue
				}

				settingName, err := valueObject.NewPhpSettingName(settingParts[0])
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}

				settingValue, err := valueObject.NewPhpSettingValue(settingParts[1])
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}

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
			runtimeCmdRepo := runtimeInfra.NewRuntimeCmdRepo(controller.persistentDbSvc)
			vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)

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

func (controller *RuntimeController) UpdatePhpSetting() *cobra.Command {
	var hostnameStr string
	var phpVersionStr string
	var settingNameStr string
	var settingValueStr string

	cmd := &cobra.Command{
		Use:   "update-setting",
		Short: "UpdatePhpSetting",
		Run: func(cmd *cobra.Command, args []string) {
			serviceName, _ := valueObject.NewServiceName("php-webserver")
			sharedHelper.StopIfServiceUnavailable(controller.persistentDbSvc, serviceName)

			hostname, err := getHostname(hostnameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			phpVersion, err := valueObject.NewPhpVersion(phpVersionStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			phpSettings := []entity.PhpSetting{}
			settingName, err := valueObject.NewPhpSettingName(settingNameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			settingValue, err := valueObject.NewPhpSettingValue(settingValueStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

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
			runtimeCmdRepo := runtimeInfra.NewRuntimeCmdRepo(controller.persistentDbSvc)
			vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)

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

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.Flags().StringVarP(&phpVersionStr, "version", "v", "", "PhpVersion")
	cmd.MarkFlagRequired("version")
	cmd.Flags().StringVarP(&settingNameStr, "name", "N", "", "PhpSettingName")
	cmd.Flags().StringVarP(&settingValueStr, "value", "V", "", "PhpSettingValue")
	return cmd
}

func (controller *RuntimeController) UpdatePhpModule() *cobra.Command {
	var hostnameStr string
	var phpVersionStr string
	var moduleNameStr string
	moduleStatusBool := true

	cmd := &cobra.Command{
		Use:   "update-module",
		Short: "UpdatePhpModule",
		Run: func(cmd *cobra.Command, args []string) {
			serviceName, _ := valueObject.NewServiceName("php-webserver")
			sharedHelper.StopIfServiceUnavailable(controller.persistentDbSvc, serviceName)

			hostname, err := getHostname(hostnameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			phpVersion, err := valueObject.NewPhpVersion(phpVersionStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			phpModules := []entity.PhpModule{}
			moduleName, err := valueObject.NewPhpModuleName(moduleNameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

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
			runtimeCmdRepo := runtimeInfra.NewRuntimeCmdRepo(controller.persistentDbSvc)
			vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)

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

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.Flags().StringVarP(&phpVersionStr, "version", "v", "", "PhpVersion")
	cmd.MarkFlagRequired("version")
	cmd.Flags().StringVarP(&moduleNameStr, "name", "N", "", "PhpModuleName")
	cmd.Flags().BoolVarP(&moduleStatusBool, "status", "V", true, "PhpModuleStatus")
	return cmd
}
