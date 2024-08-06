package apiController

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
	"github.com/speedianet/os/src/presentation/service"
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

// ReadPhpConfigs godoc
// @Summary      ReadPhpConfigs
// @Description  Get php version, modules and settings for a hostname.
// @Tags         runtime
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        hostname 	  path   string  true  "Hostname"
// @Success      200 {object} entity.PhpConfigs
// @Router       /v1/runtime/php/{hostname}/ [get]
func (controller *RuntimeController) ReadPhpConfigs(c echo.Context) error {
	requestBody := map[string]interface{}{
		"hostname": c.Param("hostname"),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.runtimeService.ReadPhpConfigs(requestBody),
	)
}

func parsePhpModules(rawPhpModules interface{}) ([]entity.PhpModule, error) {
	modules := []entity.PhpModule{}

	rawModulesSlice, assertOk := rawPhpModules.([]interface{})
	if !assertOk {
		rawModuleUnit, assertOk := rawPhpModules.(map[string]interface{})
		if !assertOk {
			return modules, errors.New("InvalidPhpModules")
		}
		rawModulesSlice = []interface{}{rawModuleUnit}
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

func parsePhpSettings(rawPhpSettings interface{}) ([]entity.PhpSetting, error) {
	settings := []entity.PhpSetting{}

	rawSettingsSlice, assertOk := rawPhpSettings.([]interface{})
	if !assertOk {
		rawSettingUnit, assertOk := rawPhpSettings.(map[string]interface{})
		if !assertOk {
			return settings, errors.New("InvalidPhpSettings")
		}
		rawPhpSettings = []interface{}{rawSettingUnit}
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

// UpdatePhpConfigs godoc
// @Summary      UpdatePhpConfigs
// @Description  Update php version, modules and settings for a hostname.
// @Tags         runtime
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        hostname 	  path   string  true  "Hostname"
// @Param        updatePhpConfigsDto	body dto.UpdatePhpConfigs	true	"modules and settings are optional."
// @Success      200 {object} object{} "PhpConfigsUpdated"
// @Router       /v1/runtime/php/{hostname}/ [put]
func (controller *RuntimeController) UpdatePhpConfigs(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}
	requestBody["hostname"] = c.Param("hostname")

	phpModules, err := parsePhpModules(requestBody["modules"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
	}
	requestBody["modules"] = phpModules

	phpSettings, err := parsePhpSettings(requestBody["settings"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
	}
	requestBody["settings"] = phpSettings

	return apiHelper.ServiceResponseWrapper(
		c, controller.runtimeService.UpdatePhpConfigs(requestBody),
	)
}
