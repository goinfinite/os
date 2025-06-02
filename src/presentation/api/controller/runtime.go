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
	"github.com/goinfinite/os/src/presentation/liaison"
	"github.com/labstack/echo/v4"
)

type RuntimeController struct {
	runtimeLiaison *liaison.RuntimeLiaison
}

func NewRuntimeController(
	persistentDbService *internalDbInfra.PersistentDatabaseService,
	trailDbService *internalDbInfra.TrailDatabaseService,
) *RuntimeController {
	return &RuntimeController{
		runtimeLiaison: liaison.NewRuntimeLiaison(persistentDbService, trailDbService),
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.runtimeLiaison.ReadPhpConfigs(requestInputData),
	)
}

func (controller *RuntimeController) parsePhpModules(rawPhpModules any) (
	[]entity.PhpModule, error,
) {
	modules := []entity.PhpModule{}

	rawModulesSlice, assertOk := rawPhpModules.([]any)
	if !assertOk {
		rawUniqueModule, assertOk := rawPhpModules.(map[string]any)
		if !assertOk {
			return modules, errors.New("InvalidPhpModulesStructure")
		}
		rawModulesSlice = []any{rawUniqueModule}
	}

	for _, rawModule := range rawModulesSlice {
		rawModuleMap, assertOk := rawModule.(map[string]any)
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

func (controller *RuntimeController) parsePhpSettings(rawPhpSettings any) (
	[]entity.PhpSetting, error,
) {
	settings := []entity.PhpSetting{}

	rawSettingsSlice, assertOk := rawPhpSettings.([]any)
	if !assertOk {
		rawUniqueSetting, assertOk := rawPhpSettings.(map[string]any)
		if !assertOk {
			return settings, errors.New("InvalidPhpSettingsStructure")
		}
		rawSettingsSlice = []any{rawUniqueSetting}
	}

	for _, rawSetting := range rawSettingsSlice {
		rawSettingMap, assertOk := rawSetting.(map[string]any)
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	if _, exists := requestInputData["modules"]; exists {
		phpModules, err := controller.parsePhpModules(requestInputData["modules"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
		}
		requestInputData["modules"] = phpModules
	}

	if _, exists := requestInputData["settings"]; exists {
		phpSettings, err := controller.parsePhpSettings(requestInputData["settings"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
		}
		requestInputData["settings"] = phpSettings
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.runtimeLiaison.UpdatePhpConfigs(requestInputData),
	)
}

// RunPhpCommand godoc
// @Summary      RunPhpCommand
// @Description  Run a php command as the webserver user for a given hostname.
// @Tags         runtime
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        runPhpCommandDto	body dto.RunPhpCommandRequest	true	"RunPhpCommandDto"
// @Success      200 {object} dto.RunPhpCommandResponse
// @Router       /v1/runtime/php/run/ [post]
func (controller *RuntimeController) RunPhpCommand(echoContext echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(echoContext)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		echoContext, controller.runtimeLiaison.RunPhpCommand(requestInputData),
	)
}
