package apiController

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
	"github.com/speedianet/os/src/presentation/service"
)

type ServicesController struct {
	servicesService *service.ServicesService
}

func NewServicesController(
	persistentDbService *internalDbInfra.PersistentDatabaseService,
) *ServicesController {
	return &ServicesController{
		servicesService: service.NewServicesService(persistentDbService),
	}
}

// ReadServices	 godoc
// @Summary      ReadServices
// @Description  List installed services and their status.
// @Tags         services
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Success      200 {array} dto.InstalledServiceWithMetrics
// @Router       /v1/services/ [get]
func (controller *ServicesController) Read(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(c, controller.servicesService.Read())
}

// ReadInstallableServices	 godoc
// @Summary      ReadInstallableServices
// @Description  List installable services.
// @Tags         services
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Success      200 {array} entity.InstallableService
// @Router       /v1/services/installables/ [get]
func (controller *ServicesController) ReadInstallables(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(
		c, controller.servicesService.ReadInstallables(),
	)
}

func parseRawPortBindings(bindings []interface{}) []string {
	rawPortBindings := []string{}
	for _, rawPortBinding := range bindings {
		rawPortBindingMap, assertOk := rawPortBinding.(map[string]interface{})
		if !assertOk {
			slog.Debug(
				"InvalidPortBindingStructure", slog.Any("portBinding", rawPortBinding),
			)
			continue
		}

		rawPortStr, err := voHelper.InterfaceToString(rawPortBindingMap["port"])
		if err != nil {
			slog.Debug(err.Error(), slog.String("port", rawPortStr))
			continue
		}

		rawProtocolStr := ""
		if _, protocolInputExists := rawPortBindingMap["protocol"]; protocolInputExists {
			rawProtocolStr, err = voHelper.InterfaceToString(rawPortBindingMap["protocol"])
			if err != nil {
				slog.Debug(err.Error(), slog.String("port", rawPortStr))
				continue
			}
		}

		rawPortBindingStr := rawPortStr + "/" + rawProtocolStr
		rawPortBindings = append(rawPortBindings, rawPortBindingStr)
	}

	return rawPortBindings
}

// CreateInstallableService godoc
// @Summary      CreateInstallableService
// @Description  Install a new installable service.
// @Tags         services
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createInstallableServiceDto	body dto.CreateInstallableService	true	"Only name is required.<br />If version is not provided, it will be 'lts'.<br />If portBindings is not provided, it wil be default service port bindings.<br />If autoCreateMapping is not provided, it will be 'true'."
// @Success      201 {object} object{} "InstallableServiceCreated"
// @Router       /v1/services/installables/ [post]
func (controller *ServicesController) CreateInstallable(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	rawEnvsSlice := []string{}
	if requestBody["envs"] != nil {
		for _, rawEnv := range requestBody["envs"].([]interface{}) {
			rawEnvStr, err := voHelper.InterfaceToString(rawEnv)
			if err != nil {
				slog.Debug(err.Error(), slog.Any("env", rawEnv))
				continue
			}
			rawEnvsSlice = append(rawEnvsSlice, rawEnvStr)
		}
	}
	requestBody["envs"] = rawEnvsSlice

	rawPortBindings := []string{}
	if requestBody["portBindings"] != nil {
		rawPortBindings = parseRawPortBindings(
			requestBody["portBindings"].([]interface{}),
		)
	}
	requestBody["portBindings"] = rawPortBindings

	return apiHelper.ServiceResponseWrapper(
		c, controller.servicesService.CreateInstallable(requestBody, true),
	)
}

// CreateCustomService godoc
// @Summary      CreateCustomService
// @Description  Install a new custom service.
// @Tags         services
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createCustomServiceDto	body dto.CreateCustomService	true	"name, type and startCmd is required.<br />If version is not provided, it will be 'lts'.<br />If portBindings is not provided, it wil be default service port bindings.<br />If autoCreateMapping is not provided, it will be 'true'."
// @Success      201 {object} object{} "CustomServiceCreated"
// @Router       /v1/services/custom/ [post]
func (controller *ServicesController) CreateCustom(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	rawEnvsSlice := []string{}
	if requestBody["envs"] != nil {
		for _, rawEnv := range requestBody["envs"].([]interface{}) {
			rawEnvStr, err := voHelper.InterfaceToString(rawEnv)
			if err != nil {
				slog.Debug(err.Error(), slog.Any("env", rawEnv))
				continue
			}
			rawEnvsSlice = append(rawEnvsSlice, rawEnvStr)
		}
	}
	requestBody["envs"] = rawEnvsSlice

	rawPortBindings := []string{}
	if requestBody["portBindings"] != nil {
		rawPortBindings = parseRawPortBindings(
			requestBody["portBindings"].([]interface{}),
		)
	}
	requestBody["portBindings"] = rawPortBindings

	return apiHelper.ServiceResponseWrapper(
		c, controller.servicesService.CreateCustom(requestBody),
	)
}

// UpdateService godoc
// @Summary      UpdateService
// @Description  Update service details.
// @Tags         services
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateServiceDto	body dto.UpdateService	true	"Only name is required.<br />Solo services can only change status.<br />status may be 'running', 'stopped' or 'uninstalled'."
// @Success      200 {object} object{} "ServiceUpdated"
// @Router       /v1/services/ [put]
func (controller *ServicesController) Update(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	rawPortBindings := []string{}
	if requestBody["portBindings"] != nil {
		rawPortBindings = parseRawPortBindings(
			requestBody["portBindings"].([]interface{}),
		)
	}
	requestBody["portBindings"] = rawPortBindings

	return apiHelper.ServiceResponseWrapper(
		c, controller.servicesService.Update(requestBody),
	)
}

// DeleteService godoc
// @Summary      DeleteService
// @Description  Delete/Uninstall a service.
// @Tags         services
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        svcName path string true "ServiceName to delete"
// @Success      200 {object} object{} "ServiceDeleted"
// @Router       /v1/services/{svcName}/ [delete]
func (controller *ServicesController) Delete(c echo.Context) error {
	requestBody := map[string]interface{}{
		"name": c.Param("svcName"),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.servicesService.Delete(requestBody),
	)
}
