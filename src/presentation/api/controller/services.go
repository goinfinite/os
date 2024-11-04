package apiController

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/goinfinite/os/src/domain/useCase"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/labstack/echo/v4"
)

type ServicesController struct {
	servicesService *service.ServicesService
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewServicesController(
	persistentDbService *internalDbInfra.PersistentDatabaseService,
) *ServicesController {
	return &ServicesController{
		servicesService: service.NewServicesService(persistentDbService),
		persistentDbSvc: persistentDbService,
	}
}

// ReadServices	 godoc
// @Summary      ReadServices
// @Description  List installed services and their status.
// @Tags         services
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id query  uint  false  "Id"
// @Param        name query  string  false  "Name"
// @Param        nature query  string  false  "Nature"
// @Param        type query  string  false  "Type"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {array} dto.ReadInstalledServicesItemsResponse
// @Router       /v1/services/ [get]
func (controller *ServicesController) ReadInstalledItems(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.servicesService.ReadInstalledItems(requestBody),
	)
}

// ReadInstallableServices	 godoc
// @Summary      ReadInstallableServices
// @Description  List installable services.
// @Tags         services
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id query  uint  false  "Id"
// @Param        name query  string  false  "Name"
// @Param        nature query  string  false  "Nature"
// @Param        type query  string  false  "Type"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {array} dto.ReadInstallableServicesItemsResponse
// @Router       /v1/services/installables/ [get]
func (controller *ServicesController) ReadInstallables(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.servicesService.ReadInstallableItems(requestBody),
	)
}

func (controller *ServicesController) parseRawEnvs(envs interface{}) ([]string, error) {
	rawEnvs := []string{}
	rawEnvsSlice, assertOk := envs.([]interface{})
	if !assertOk {
		rawEnvUnique, assertOk := envs.(string)
		if !assertOk {
			return rawEnvs, errors.New("EnvsMustBeStringOrStringSlice")
		}
		rawEnvsSlice = []interface{}{rawEnvUnique}
	}

	for _, rawEnv := range rawEnvsSlice {
		rawEnvStr, err := voHelper.InterfaceToString(rawEnv)
		if err != nil {
			slog.Debug(err.Error(), slog.Any("env", rawEnv))
			continue
		}
		rawEnvs = append(rawEnvs, rawEnvStr)
	}

	return rawEnvs, nil
}

func (controller *ServicesController) parseRawPortBindings(
	bindings interface{},
) ([]string, error) {
	rawPortBindings := []string{}
	rawPortBindingsSlice, assertOk := bindings.([]interface{})
	if !assertOk {
		rawPortBindingUnique, assertOk := bindings.(map[string]interface{})
		if !assertOk {
			return rawPortBindings, errors.New("PortBindingsMustBeMapOrMapSlice")
		}
		rawPortBindingsSlice = []interface{}{rawPortBindingUnique}
	}

	for _, rawPortBinding := range rawPortBindingsSlice {
		rawPortBindingMap, assertOk := rawPortBinding.(map[string]interface{})
		if !assertOk {
			slog.Debug(
				"InvalidPortBindingStructure", slog.Any("portBinding", rawPortBinding),
			)
			continue
		}

		rawPortStr, err := voHelper.InterfaceToString(rawPortBindingMap["port"])
		if err != nil {
			slog.Debug(err.Error(), slog.Any("port", rawPortBindingMap["port"]))
			continue
		}

		rawPortBindingStr := rawPortStr

		if _, protocolInputExists := rawPortBindingMap["protocol"]; protocolInputExists {
			rawProtocolStr, err := voHelper.InterfaceToString(
				rawPortBindingMap["protocol"],
			)
			if err != nil {
				slog.Debug(err.Error(), slog.Any(
					"protocol", rawPortBindingMap["protocol"]),
				)
				continue
			}
			rawPortBindingStr += "/" + rawProtocolStr
		}

		rawPortBindings = append(rawPortBindings, rawPortBindingStr)
	}

	return rawPortBindings, nil
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

	rawEnvs := []string{}
	if requestBody["envs"] != nil {
		rawEnvs, err = controller.parseRawEnvs(requestBody["envs"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
	}
	requestBody["envs"] = rawEnvs

	rawPortBindings := []string{}
	if requestBody["portBindings"] != nil {
		rawPortBindings, err = controller.parseRawPortBindings(
			requestBody["portBindings"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
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

	rawEnvs := []string{}
	if requestBody["envs"] != nil {
		rawEnvs, err = controller.parseRawEnvs(requestBody["envs"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
	}
	requestBody["envs"] = rawEnvs

	rawPortBindings := []string{}
	if requestBody["portBindings"] != nil {
		rawPortBindings, err = controller.parseRawPortBindings(
			requestBody["portBindings"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
		}
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
// @Param        updateServiceDto	body dto.UpdateService	true	"Only name is required.<br />Solo services can only change status.<br />status may be 'running', 'stopped', 'uninstalled' or 'restarting'."
// @Success      200 {object} object{} "ServiceUpdated"
// @Router       /v1/services/ [put]
func (controller *ServicesController) Update(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	if requestBody["envs"] != nil {
		rawEnvs, err := controller.parseRawEnvs(requestBody["envs"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		requestBody["envs"] = rawEnvs
	}

	if requestBody["portBindings"] != nil {
		rawPortBindings, err := controller.parseRawPortBindings(
			requestBody["portBindings"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
		}
		requestBody["portBindings"] = rawPortBindings
	}

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

func (controller *ServicesController) AutoRefreshServicesItems() {
	taskInterval := time.Duration(2) * time.Minute
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(
		controller.persistentDbSvc,
	)
	for range timer.C {
		useCase.RefreshServicesItems(servicesCmdRepo)
	}
}
