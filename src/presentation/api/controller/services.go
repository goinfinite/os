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
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *ServicesController {
	return &ServicesController{
		servicesService: service.NewServicesService(persistentDbService, trailDbSvc),
		persistentDbSvc: persistentDbService,
	}
}

// ReadInstalledItems	 godoc
// @Summary      ReadInstalledItems
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
// @Success      200 {object} dto.ReadInstalledServicesItemsResponse
// @Router       /v1/services/ [get]
func (controller *ServicesController) ReadInstalledItems(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.servicesService.ReadInstalledItems(requestInputData),
	)
}

// ReadInstallableItems	 godoc
// @Summary      ReadInstallableItems
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
// @Success      200 {object} dto.ReadInstallableServicesItemsResponse
// @Router       /v1/services/installables/ [get]
func (controller *ServicesController) ReadInstallablesItems(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.servicesService.ReadInstallableItems(requestInputData),
	)
}

func (controller *ServicesController) parseRawEnvs(envsUnknownType any) ([]string, error) {
	rawEnvsStringSlice, assertOk := envsUnknownType.([]string)
	if assertOk {
		return rawEnvsStringSlice, nil
	}

	rawEnvsMap := map[string]interface{}{}
	switch envsValues := envsUnknownType.(type) {
	case []string:
		return envsValues, nil
	case string:
		return []string{envsValues}, nil
	case []interface{}:
		for _, envInterface := range envsValues {
			envMap, assertOk := envInterface.(map[string]interface{})
			if !assertOk {
				slog.Debug("InvalidEnvStructure", slog.Any("envVar", envMap["key"]))
				continue
			}

			envMapKeyNameStr, err := voHelper.InterfaceToString(envMap["key"])
			if err != nil {
				slog.Debug(err.Error(), slog.Any("envVar", envMap["key"]))
				continue
			}
			rawEnvsMap[envMapKeyNameStr] = envMap["value"]
		}
	case map[string]interface{}:
		rawEnvsMap = envsValues
	default:
		return []string{}, errors.New("EnvsMustBeStringOrStringSliceOrMapOrMapSlice")
	}

	rawEnvsStrSlice := []string{}
	for mapPropName, mapPropValue := range rawEnvsMap {
		mapPropValueStr, err := voHelper.InterfaceToString(mapPropValue)
		if err != nil {
			slog.Debug("InvalidEnvValue", slog.Any("envVar", mapPropName))
			continue
		}

		rawEnvsStrSlice = append(rawEnvsStrSlice, mapPropName+"="+mapPropValueStr)
	}

	return rawEnvsStrSlice, nil
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	rawEnvs := []string{}
	if requestInputData["envs"] != nil {
		rawEnvs, err = controller.parseRawEnvs(requestInputData["envs"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
	}
	requestInputData["envs"] = rawEnvs

	rawPortBindings := []string{}
	if requestInputData["portBindings"] != nil {
		rawPortBindings, err = controller.parseRawPortBindings(
			requestInputData["portBindings"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
	}
	requestInputData["portBindings"] = rawPortBindings

	return apiHelper.ServiceResponseWrapper(
		c, controller.servicesService.CreateInstallable(requestInputData, true),
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	rawEnvs := []string{}
	if requestInputData["envs"] != nil {
		rawEnvs, err = controller.parseRawEnvs(requestInputData["envs"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
	}
	requestInputData["envs"] = rawEnvs

	rawPortBindings := []string{}
	if requestInputData["portBindings"] != nil {
		rawPortBindings, err = controller.parseRawPortBindings(
			requestInputData["portBindings"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
		}
	}
	requestInputData["portBindings"] = rawPortBindings

	return apiHelper.ServiceResponseWrapper(
		c, controller.servicesService.CreateCustom(requestInputData),
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	if requestInputData["envs"] != nil {
		rawEnvs, err := controller.parseRawEnvs(requestInputData["envs"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		requestInputData["envs"] = rawEnvs
	}

	if requestInputData["portBindings"] != nil {
		rawPortBindings, err := controller.parseRawPortBindings(
			requestInputData["portBindings"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
		}
		requestInputData["portBindings"] = rawPortBindings
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.servicesService.Update(requestInputData),
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.servicesService.Delete(requestInputData),
	)
}

func (controller *ServicesController) AutoRefreshServiceInstallableItems() {
	refreshIntervalHours := 24 / useCase.RefreshServiceInstallableItemsAmountPerDay

	taskInterval := time.Duration(refreshIntervalHours) * time.Hour
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(
		controller.persistentDbSvc,
	)
	for range timer.C {
		useCase.RefreshServiceInstallableItems(servicesCmdRepo)
	}
}
