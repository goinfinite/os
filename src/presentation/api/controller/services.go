package apiController

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
	sharedHelper "github.com/goinfinite/os/src/presentation/shared/helper"
	"github.com/labstack/echo/v4"
)

type ServicesController struct {
	servicesLiaison *liaison.ServicesLiaison
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewServicesController(
	persistentDbService *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *ServicesController {
	return &ServicesController{
		servicesLiaison: liaison.NewServicesLiaison(persistentDbService, trailDbSvc),
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

	return apiHelper.LiaisonResponseWrapper(
		c, controller.servicesLiaison.ReadInstalledItems(requestInputData),
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

	return apiHelper.LiaisonResponseWrapper(
		c, controller.servicesLiaison.ReadInstallableItems(requestInputData),
	)
}

func (controller *ServicesController) transformRawEnvsInterfaceSliceToMap(
	rawEnvsInterface []interface{},
) map[string]interface{} {
	rawEnvsMap := map[string]interface{}{}

	for rawEnvIndex, rawEnvInterface := range rawEnvsInterface {
		rawEnvMap, assertOk := rawEnvInterface.(map[string]interface{})
		if !assertOk {
			slog.Debug("InvalidEnvStructure", slog.Int("envIndex", rawEnvIndex))
			continue
		}

		envVarNameStr, err := voHelper.InterfaceToString(rawEnvMap["name"])
		if err != nil {
			slog.Debug("EnvVarNameMustBeString", slog.Int("envIndex", rawEnvIndex))
			continue
		}
		rawEnvsMap[envVarNameStr] = rawEnvMap["value"]
	}

	return rawEnvsMap
}

func (controller *ServicesController) parseRawEnvs(
	rawEnvsUnknownType any,
) ([]valueObject.ServiceEnv, error) {
	serviceEnvs := []valueObject.ServiceEnv{}

	rawEnvsMap := map[string]interface{}{}
	switch rawEnvsValues := rawEnvsUnknownType.(type) {
	case string, []string:
		return sharedHelper.StringSliceValueObjectParser(
			rawEnvsValues, valueObject.NewServiceEnv,
		), nil
	case []interface{}:
		rawEnvsMap = controller.transformRawEnvsInterfaceSliceToMap(rawEnvsValues)
	case map[string]interface{}:
		rawEnvsMap = rawEnvsValues
	default:
		return serviceEnvs, errors.New("EnvsMustBeStringOrStringSliceOrMapOrMapSlice")
	}

	for envVarName, envVarValue := range rawEnvsMap {
		envValueStr, err := voHelper.InterfaceToString(envVarValue)
		if err != nil {
			slog.Debug("InvalidServiceEnvValue", slog.Any("envVarName", envVarName))
			continue
		}

		serviceEnvStr := envVarName + "=" + envValueStr
		serviceEnv, err := valueObject.NewServiceEnv(serviceEnvStr)
		if err != nil {
			slog.Debug(err.Error(), slog.Any("envVarName", envVarName))
			continue
		}

		serviceEnvs = append(serviceEnvs, serviceEnv)
	}

	return serviceEnvs, nil
}

func (controller *ServicesController) transformRawPortBindingsInterfaceSliceToMap(
	rawPortBindingsInterface []interface{},
) map[string]interface{} {
	rawPortBindingsMap := map[string]interface{}{}

	for rawPortBindingIndex, rawPortBindingInterface := range rawPortBindingsInterface {
		rawPortBindingMap, assertOk := rawPortBindingInterface.(map[string]interface{})
		if !assertOk {
			slog.Debug(
				"InvalidPortBindingStructure",
				slog.Int("portBindingIndex", rawPortBindingIndex),
			)
			continue
		}

		portStr, err := voHelper.InterfaceToString(rawPortBindingMap["port"])
		if err != nil {
			slog.Debug(
				"PortBindingPortBeString",
				slog.Int("portBindingIndex", rawPortBindingIndex),
			)
			continue
		}
		rawPortBindingsMap[portStr] = rawPortBindingMap["protocol"]
	}

	return rawPortBindingsMap
}

func (controller *ServicesController) parseRawPortBindings(
	rawPortBindingsUnknownType any,
) ([]valueObject.PortBinding, error) {
	portBindings := []valueObject.PortBinding{}

	rawPortBindingsInterfaceSlice := []interface{}{}
	switch rawPortBindingsValues := rawPortBindingsUnknownType.(type) {
	case string, []string:
		return sharedHelper.StringSliceValueObjectParser(
			rawPortBindingsValues, valueObject.NewPortBinding,
		), nil
	case []interface{}:
		rawPortBindingsInterfaceSlice = rawPortBindingsValues
	case map[string]interface{}:
		rawPortBindingsInterfaceSlice = []interface{}{rawPortBindingsValues}
	default:
		return portBindings, errors.New(
			"PortBindingsMustBeStringOrStringSliceOrMapOrMapSlice",
		)
	}

	rawPortBindingsMap := controller.transformRawPortBindingsInterfaceSliceToMap(
		rawPortBindingsInterfaceSlice,
	)

	rawPortBindingIndex := -1
	for port, protocol := range rawPortBindingsMap {
		rawPortBindingIndex++

		protocolStr, err := voHelper.InterfaceToString(protocol)
		if err != nil {
			slog.Debug(
				"InvalidPortBindingProtocol",
				slog.Int("portBindingIndex", rawPortBindingIndex),
			)
			continue
		}

		portBindingStr := port + "/" + protocolStr
		portBinding, err := valueObject.NewPortBinding(portBindingStr)
		if err != nil {
			slog.Debug(err.Error(), slog.Any("port", port))
			continue
		}

		portBindings = append(portBindings, portBinding)
	}

	return portBindings, nil
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

	envs := []valueObject.ServiceEnv{}
	if requestInputData["envs"] != nil {
		envs, err = controller.parseRawEnvs(requestInputData["envs"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
	}
	requestInputData["envs"] = envs

	portBindings := []valueObject.PortBinding{}
	if requestInputData["portBindings"] != nil {
		portBindings, err = controller.parseRawPortBindings(
			requestInputData["portBindings"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
	}
	requestInputData["portBindings"] = portBindings

	return apiHelper.LiaisonResponseWrapper(
		c, controller.servicesLiaison.CreateInstallable(requestInputData, true),
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

	envs := []valueObject.ServiceEnv{}
	if requestInputData["envs"] != nil {
		envs, err = controller.parseRawEnvs(requestInputData["envs"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
	}
	requestInputData["envs"] = envs

	portBindings := []valueObject.PortBinding{}
	if requestInputData["portBindings"] != nil {
		portBindings, err = controller.parseRawPortBindings(
			requestInputData["portBindings"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
		}
	}
	requestInputData["portBindings"] = portBindings

	return apiHelper.LiaisonResponseWrapper(
		c, controller.servicesLiaison.CreateCustom(requestInputData),
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

	return apiHelper.LiaisonResponseWrapper(
		c, controller.servicesLiaison.Update(requestInputData),
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
	requestInputData["name"] = requestInputData["svcName"]

	return apiHelper.LiaisonResponseWrapper(
		c, controller.servicesLiaison.Delete(requestInputData),
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
