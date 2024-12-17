package apiController

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/labstack/echo/v4"
)

type RuntimeController struct {
	runtimeService *service.RuntimeService
}

func NewRuntimeController(
	persistentDbService *internalDbInfra.PersistentDatabaseService,
	trailDbService *internalDbInfra.TrailDatabaseService,
) *RuntimeController {
	return &RuntimeController{
		runtimeService: service.NewRuntimeService(persistentDbService, trailDbService),
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

func (controller *RuntimeController) parsePhpModules(rawPhpModules interface{}) (
	[]entity.PhpModule, error,
) {
	modules := []entity.PhpModule{}

	rawModulesSlice, assertOk := rawPhpModules.([]interface{})
	if !assertOk {
		rawUniqueModule, assertOk := rawPhpModules.(map[string]interface{})
		if !assertOk {
			return modules, errors.New("InvalidPhpModulesStructure")
		}
		rawModulesSlice = []interface{}{rawUniqueModule}
	}

	for _, rawModule := range rawModulesSlice {
		rawModuleMap, assertOk := rawModule.(map[string]interface{})
		if !assertOk {
			slog.Debug("PhpModuleIsNotAnInterface")
			continue
		}

		moduleName, err := valueObject.NewPhpModuleName(rawModuleMap["name"])
		if err != nil {
			slog.Debug(err.Error(), slog.Any("name", rawModuleMap["name"]))
			continue
		}

		moduleStatus, err := voHelper.InterfaceToBool(rawModuleMap["status"])
		if err != nil {
			slog.Debug(err.Error(), slog.Any("status", rawModuleMap["status"]))
			continue
		}

		modules = append(modules, entity.NewPhpModule(moduleName, moduleStatus))
	}

	return modules, nil
}

func (controller *RuntimeController) parsePhpSettings(rawPhpSettings interface{}) (
	[]entity.PhpSetting, error,
) {
	settings := []entity.PhpSetting{}

	rawSettingsSlice, assertOk := rawPhpSettings.([]interface{})
	if !assertOk {
		rawUniqueSetting, assertOk := rawPhpSettings.(map[string]interface{})
		if !assertOk {
			return settings, errors.New("InvalidPhpSettingsStructure")
		}
		rawSettingsSlice = []interface{}{rawUniqueSetting}
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

		settingType, _ := valueObject.NewPhpSettingType("text")

		settings = append(
			settings,
			entity.NewPhpSetting(
				settingName, settingType, settingValue, []valueObject.PhpSettingOption{},
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
	requestBody, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}
	requestBody["hostname"] = c.Param("hostname")

	if _, exists := requestBody["modules"]; exists {
		phpModules, err := controller.parsePhpModules(requestBody["modules"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
		}
		requestBody["modules"] = phpModules
	}

	if _, exists := requestBody["settings"]; exists {
		phpSettings, err := controller.parsePhpSettings(requestBody["settings"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
		}
		requestBody["settings"] = phpSettings
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.runtimeService.UpdatePhpConfigs(requestBody),
	)
}
