package apiController

import (
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

func parsePhpModules(
	rawPhpModules map[string]interface{},
) ([]entity.PhpModule, error) {
	modules := []entity.PhpModule{}
	rawModules, ok := rawPhpModules["modules"].([]interface{})
	if !ok {
		return modules, nil
	}

	for _, rawModule := range rawModules {
		rawModuleMap, ok := rawModule.(map[string]interface{})
		if !ok {
			continue
		}

		moduleName, err := valueObject.NewPhpModuleName(rawModuleMap["name"])
		if err != nil {
			continue
		}

		moduleStatus, ok := rawModuleMap["status"].(bool)
		if !ok {
			continue
		}

		modules = append(modules, entity.NewPhpModule(moduleName, moduleStatus))
	}

	return modules, nil
}

func parsePhpSettings(
	rawPhpSettings map[string]interface{},
) ([]entity.PhpSetting, error) {
	settings := []entity.PhpSetting{}
	rawSettings, ok := rawPhpSettings["settings"].([]interface{})
	if !ok {
		return settings, nil
	}

	for _, rawSetting := range rawSettings {
		rawSettingMap, ok := rawSetting.(map[string]interface{})
		if !ok {
			continue
		}

		settingName, err := valueObject.NewPhpSettingName(rawSettingMap["name"])
		if err != nil {
			continue
		}

		settingValue, err := valueObject.NewPhpSettingValue(rawSettingMap["value"])
		if err != nil {
			continue
		}

		emptySettingOptions := []valueObject.PhpSettingOption{}

		settings = append(
			settings,
			entity.NewPhpSetting(settingName, settingValue, emptySettingOptions),
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

	phpModules, err := parsePhpModules(requestBody)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}
	requestBody["modules"] = phpModules

	phpSettings, err := parsePhpSettings(requestBody)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}
	requestBody["settings"] = phpSettings

	return apiHelper.ServiceResponseWrapper(
		c, controller.runtimeService.UpdatePhpConfigs(requestBody),
	)
}
