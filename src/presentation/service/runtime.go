package service

import (
	"errors"
	"log/slog"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	runtimeInfra "github.com/speedianet/os/src/infra/runtime"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	serviceHelper "github.com/speedianet/os/src/presentation/service/helper"
	sharedHelper "github.com/speedianet/os/src/presentation/shared/helper"
)

type RuntimeService struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewRuntimeService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *RuntimeService {
	return &RuntimeService{
		persistentDbSvc: persistentDbSvc,
	}
}

func (service *RuntimeService) ReadPhpConfigs(
	input map[string]interface{},
) ServiceOutput {
	serviceName, _ := valueObject.NewServiceName("php-webserver")
	sharedHelper.StopIfServiceUnavailable(service.persistentDbSvc, serviceName)

	hostname, err := valueObject.NewFqdn(input["hostname"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	runtimeQueryRepo := runtimeInfra.RuntimeQueryRepo{}
	phpConfigs, err := useCase.ReadPhpConfigs(runtimeQueryRepo, hostname)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, phpConfigs)
}

func (service *RuntimeService) parsePhpModules(
	rawPhpModules interface{},
) ([]entity.PhpModule, error) {
	modules := []entity.PhpModule{}

	rawModulesSlice, assertOk := rawPhpModules.([]interface{})
	if !assertOk {
		return modules, errors.New("InvalidPhpModules")
	}

	for _, rawModule := range rawModulesSlice {
		rawModuleMap, assertOk := rawModule.(map[string]interface{})
		if !assertOk {
			slog.Debug("PhpModuleIsNotAnInterface")
			continue
		}

		moduleName, err := valueObject.NewPhpModuleName(rawModuleMap["name"])
		if err != nil {
			slog.Debug("InvalidPhpModuleName", slog.Any("name", rawModuleMap["name"]))
			continue
		}

		moduleStatus, assertOk := rawModuleMap["status"].(bool)
		if !assertOk {
			slog.Debug(
				"InvalidPhpModuleStatus", slog.Any("status", rawModuleMap["status"]),
			)
			continue
		}

		modules = append(modules, entity.NewPhpModule(moduleName, moduleStatus))
	}

	return modules, nil
}

func (service *RuntimeService) parsePhpSettings(
	rawPhpSettings interface{},
) ([]entity.PhpSetting, error) {
	settings := []entity.PhpSetting{}

	rawSettingsSlice, assertOk := rawPhpSettings.([]interface{})
	if !assertOk {
		return settings, errors.New("InvalidPhpSettings")
	}

	for _, rawSetting := range rawSettingsSlice {
		rawSettingMap, assertOk := rawSetting.(map[string]interface{})
		if !assertOk {
			slog.Debug("PhpSettingIsNotAnInterface")
			continue
		}

		settingName, err := valueObject.NewPhpSettingName(rawSettingMap["name"])
		if err != nil {
			slog.Debug(
				"InvalidPhpSettingName", slog.Any("name", rawSettingMap["name"]),
			)
			continue
		}

		settingValue, err := valueObject.NewPhpSettingValue(rawSettingMap["value"])
		if err != nil {
			slog.Debug(
				"InvalidPhpSettingValue", slog.Any("value", rawSettingMap["value"]),
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

	return settings, nil
}

func (service *RuntimeService) UpdatePhpConfigs(
	input map[string]interface{},
) ServiceOutput {
	serviceName, _ := valueObject.NewServiceName("php-webserver")
	sharedHelper.StopIfServiceUnavailable(service.persistentDbSvc, serviceName)

	requiredParams := []string{"hostname", "version"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	hostname, err := valueObject.NewFqdn(input["hostname"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	phpVersion, err := valueObject.NewPhpVersion(input["version"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	phpModules := []entity.PhpModule{}
	if input["modules"] != nil {
		phpModules, err = service.parsePhpModules(input["modules"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	phpSettings := []entity.PhpSetting{}
	if input["settings"] != nil {
		phpSettings, err = service.parsePhpSettings(input["settings"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	dto := dto.NewUpdatePhpConfigs(hostname, phpVersion, phpModules, phpSettings)

	runtimeQueryRepo := runtimeInfra.RuntimeQueryRepo{}
	runtimeCmdRepo := runtimeInfra.NewRuntimeCmdRepo(service.persistentDbSvc)
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbSvc)

	err = useCase.UpdatePhpConfigs(
		runtimeQueryRepo, runtimeCmdRepo, vhostQueryRepo, dto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "PhpConfigsUpdated")
}
