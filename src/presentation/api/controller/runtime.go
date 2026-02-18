package apiController

import (
	tkPresentation "github.com/goinfinite/tk/src/presentation"
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
func (controller *RuntimeController) ReadPhpConfigs(echoContext echo.Context) error {
	inputReader := tkPresentation.ApiRequestInputReader{}
	requestData, requestParsingErr := inputReader.Reader(echoContext)
	if requestParsingErr != nil {
		return requestParsingErr
	}

	return apiHelper.LiaisonResponseWrapper(
		echoContext, controller.runtimeLiaison.ReadPhpConfigs(requestData),
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
func (controller *RuntimeController) UpdatePhpConfigs(echoContext echo.Context) error {
	inputReader := tkPresentation.ApiRequestInputReader{}
	requestData, requestParsingErr := inputReader.Reader(echoContext)
	if requestParsingErr != nil {
		return requestParsingErr
	}

	if _, exists := requestData["modules"]; exists {
		phpModules, err := controller.parsePhpModules(requestData["modules"])
		if err != nil {
			return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err)
		}
		requestData["modules"] = phpModules
	}

	if _, exists := requestData["settings"]; exists {
		phpSettings, err := controller.parsePhpSettings(requestData["settings"])
		if err != nil {
			return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err)
		}
		requestData["settings"] = phpSettings
	}

	return apiHelper.LiaisonResponseWrapper(
		echoContext, controller.runtimeLiaison.UpdatePhpConfigs(requestData),
	)
}

// RunPhpCommand godoc
// @Summary      RunPhpCommand
// @Description  Run a php command as the webserver user for a given hostname. <br />CAUTION: This endpoint allows for arbitrary code execution (ACE) and is therefore disabled by default. <br />To enable this endpoint, set the "ENABLE_API_RUNTIME_PHP_RUN_CMD" environment variable to "true" when starting the API/container.<br />Only super admin accounts can use this endpoint.
// @Tags         runtime
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        runPhpCmdDto	body dto.RunPhpCommandRequest	true	"Hostname and command are required. Timeout is optional."
// @Success      200 {object} dto.RunPhpCommandResponse
// @Router       /v1/runtime/php/run/ [post]
func (controller *RuntimeController) RunPhpCommand(echoContext echo.Context) error {
	requestData, err := apiHelper.ReadRequestInputData(echoContext)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		echoContext, controller.runtimeLiaison.RunPhpCommand(requestData),
	)
}
