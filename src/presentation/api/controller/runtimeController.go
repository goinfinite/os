package apiController

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	runtimeInfra "github.com/speedianet/os/src/infra/runtime"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
	sharedHelper "github.com/speedianet/os/src/presentation/shared/helper"
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
	serviceName, _ := valueObject.NewServiceName("php-webserver")
	sharedHelper.StopIfServiceUnavailable(controller.persistentDbSvc, serviceName)

	hostname, err := valueObject.NewFqdn(c.Param("hostname"))
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	runtimeQueryRepo := runtimeInfra.RuntimeQueryRepo{}
	phpConfigs, err := useCase.ReadPhpConfigs(runtimeQueryRepo, hostname)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, phpConfigs)
}

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

		moduleName, err := valueObject.NewPhpModuleName(moduleMap["name"])
		if err != nil {
			return nil, err
		}

		moduleStatus, ok := moduleMap["status"].(bool)
		if !ok {
			return nil, errors.New("InvalidModuleStatus")
		}

		phpModules = append(phpModules, entity.NewPhpModule(moduleName, moduleStatus))
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

		settingName, err := valueObject.NewPhpSettingName(settingMap["name"])
		if err != nil {
			return nil, err
		}

		settingValue, err := valueObject.NewPhpSettingValue(settingMap["value"])
		if err != nil {
			return nil, err
		}

		phpSettings = append(
			phpSettings,
			entity.NewPhpSetting(
				settingName, settingValue, []valueObject.PhpSettingOption{},
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
// @Param        hostname 	  path   string  true  "Hostname"
// @Param        updatePhpConfigsDto	body dto.UpdatePhpConfigs	true	"modules and settings are optional."
// @Success      200 {object} object{} "PhpConfigsUpdated"
// @Router       /v1/runtime/php/{hostname}/ [put]
func (controller *RuntimeController) UpdatePhpConfigs(c echo.Context) error {
	serviceName, _ := valueObject.NewServiceName("php-webserver")
	sharedHelper.StopIfServiceUnavailable(controller.persistentDbSvc, serviceName)

	hostname, err := valueObject.NewFqdn(c.Param("hostname"))
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	requiredParams := []string{"version"}
	requestBody, _ := apiHelper.ReadRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	phpVersion, err := valueObject.NewPhpVersion(requestBody["version"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

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
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "PhpConfigsUpdated")
}
