package apiController

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/domain/valueObject"
	"github.com/speedianet/sam/src/infra"
	apiHelper "github.com/speedianet/sam/src/presentation/api/helper"
)

func getPhpModules(requestBody map[string]interface{}) ([]entity.PhpModule, error) {
	var phpModules []entity.PhpModule
	modules, ok := requestBody["modules"].([]interface{})
	if !ok {
		return nil, nil
	}

	for _, module := range modules {
		moduleMap, ok := module.(map[string]interface{})
		if !ok {
			return nil, errors.New("InvalidModuleStruct")
		}

		moduleName, ok := moduleMap["name"].(string)
		if !ok {
			return nil, errors.New("InvalidModuleName")
		}

		moduleStatus, ok := moduleMap["status"].(bool)
		if !ok {
			return nil, errors.New("InvalidModuleStatus")
		}

		phpModules = append(
			phpModules,
			entity.NewPhpModule(
				valueObject.NewPhpModuleNamePanic(moduleName),
				moduleStatus,
			),
		)
	}

	return phpModules, nil
}

func getPhpSettings(requestBody map[string]interface{}) ([]entity.PhpSetting, error) {
	var phpSettings []entity.PhpSetting
	settings, ok := requestBody["settings"].([]interface{})
	if !ok {
		return nil, nil
	}

	for _, setting := range settings {
		settingMap, ok := setting.(map[string]interface{})
		if !ok {
			return nil, errors.New("InvalidSettingStruct")
		}

		settingName, ok := settingMap["name"].(string)
		if !ok {
			return nil, errors.New("InvalidSettingName")
		}

		valueSent := settingMap["value"]
		var settingValue string
		switch value := valueSent.(type) {
		case string:
			settingValue = value
		case bool:
			settingValue = strconv.FormatBool(value)
		case int:
			settingValue = strconv.Itoa(value)
		case float64:
			settingValue = strconv.FormatFloat(value, 'f', -1, 64)
		default:
			return nil, errors.New("InvalidSettingValue")
		}

		phpSettings = append(
			phpSettings,
			entity.NewPhpSetting(
				valueObject.NewPhpSettingNamePanic(settingName),
				valueObject.NewPhpSettingValuePanic(settingValue),
				[]valueObject.PhpSettingOption{},
			),
		)
	}

	return phpSettings, nil
}

// UpdatePhpConfigs godoc
// @Summary      UpdatePhpConfigs
// @Description  Update php version, modules and settings for a hostname.
// @Tags         runtime
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updatePhpConfigsDto	body dto.UpdatePhpConfigs	true	"UpdatePhpConfigs"
// @Success      200 {object} object{} "PhpConfigsUpdated"
// @Router       /runtime/php/{hostname}/ [put]
func UpdatePhpConfigsController(c echo.Context) error {
	hostname := valueObject.NewFqdnPanic(c.Param("hostname"))

	requiredParams := []string{"version"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	phpVersion := valueObject.NewPhpVersionPanic(requestBody["version"].(string))

	phpModules, err := getPhpModules(requestBody)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	phpSettings, err := getPhpSettings(requestBody)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	updatePhpConfigsDto := dto.NewUpdatePhpConfigs(
		hostname,
		phpVersion,
		phpModules,
		phpSettings,
	)

	runtimeQueryRepo := infra.RuntimeQueryRepo{}
	runtimeCmdRepo := infra.RuntimeCmdRepo{}

	err = useCase.UpdatePhpConfigs(
		runtimeQueryRepo,
		runtimeCmdRepo,
		updatePhpConfigsDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "PhpConfigsUpdated")
}
