package apiController

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
	"github.com/speedianet/os/src/presentation/service"
)

type ServicesController struct {
	serviceService *service.ServicesService
}

func NewServicesController(
	persistentDbService *internalDbInfra.PersistentDatabaseService,
) *ServicesController {
	return &ServicesController{
		serviceService: service.NewServicesService(persistentDbService),
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
	return apiHelper.ServiceResponseWrapper(c, controller.serviceService.Read())
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
		c, controller.serviceService.ReadInstallables(),
	)
}

func parseRawPortBindings(bindings []interface{}) []string {
	rawPortBindings := []string{}
	for _, rawPortBinding := range bindings {
		rawPortBindingMap, assertOk := rawPortBinding.(map[string]interface{})
		if !assertOk {
			continue
		}

		rawPortStr, assertOk := rawPortBindingMap["port"].(string)
		if !assertOk {
			rawPortFloat, assertOk := rawPortBindingMap["port"].(float64)
			if !assertOk {
				continue
			}
			rawPortStr = strconv.FormatFloat(rawPortFloat, 'f', -1, 64)
		}

		rawProtocolStr, assertOk := rawPortBindingMap["protocol"].(string)
		if !assertOk {
			continue
		}

		rawPortBinding = rawPortStr + "/" + rawProtocolStr
		rawPortBindings = append(rawPortBindings, rawPortBinding.(string))
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
			rawEnvsSlice = append(rawEnvsSlice, rawEnv.(string))
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

	if requestBody["timeoutStartSecs"] != nil {
		requestBody["timeoutStartSecs"], err = voHelper.InterfaceToUint(
			requestBody["timeoutStartSecs"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "TimeoutStartSecsMustBeUint",
			)
		}
	}

	if requestBody["maxStartRetries"] != nil {
		requestBody["maxStartRetries"], err = voHelper.InterfaceToUint(
			requestBody["maxStartRetries"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "MaxStartRetriesMustBeUint",
			)
		}
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.serviceService.CreateInstallable(requestBody, true),
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
			rawEnvsSlice = append(rawEnvsSlice, rawEnv.(string))
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

	if requestBody["timeoutStartSecs"] != nil {
		requestBody["timeoutStartSecs"], err = voHelper.InterfaceToUint(
			requestBody["timeoutStartSecs"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "TimeoutStartSecsMustBeUint",
			)
		}
	}

	if requestBody["maxStartRetries"] != nil {
		requestBody["maxStartRetries"], err = voHelper.InterfaceToUint(
			requestBody["maxStartRetries"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "MaxStartRetriesMustBeUint",
			)
		}
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.serviceService.CreateCustom(requestBody),
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
		c, controller.serviceService.Update(requestBody),
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
		c, controller.serviceService.Delete(requestBody),
	)
}
