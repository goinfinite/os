package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	mappingInfra "github.com/speedianet/os/src/infra/vhost/mapping"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
	sharedHelper "github.com/speedianet/os/src/presentation/shared/helper"
)

type ServicesController struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewServicesController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ServicesController {
	return &ServicesController{
		persistentDbSvc: persistentDbSvc,
	}
}

// ReadServices	 godoc
// @Summary      ReadServices
// @Description  List installed services and their status.
// @Tags         services
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Success      200 {array} dto.ServiceWithMetrics
// @Router       /v1/services/ [get]
func (controller *ServicesController) Read(c echo.Context) error {
	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	servicesList, err := useCase.GetServicesWithMetrics(servicesQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, servicesList)
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
	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	servicesList, err := useCase.GetInstallableServices(servicesQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, servicesList)
}

func parsePortBindings(bindings []interface{}) []valueObject.PortBinding {
	var svcPortBindings []valueObject.PortBinding
	for _, portBinding := range bindings {
		portBindingMap := portBinding.(map[string]interface{})
		svcPort := valueObject.NewNetworkPortPanic(
			portBindingMap["port"],
		)
		svcProtocol := valueObject.NewNetworkProtocolPanic(
			portBindingMap["protocol"].(string),
		)
		svcPortBinding := valueObject.NewPortBinding(
			svcPort,
			svcProtocol,
		)
		svcPortBindings = append(svcPortBindings, svcPortBinding)
	}
	return svcPortBindings
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
	requiredParams := []string{"name"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	svcName := valueObject.NewServiceNamePanic(requestBody["name"].(string))

	var svcVersionPtr *valueObject.ServiceVersion
	if requestBody["version"] != nil {
		svcVersion := valueObject.NewServiceVersionPanic(
			requestBody["version"].(string),
		)
		svcVersionPtr = &svcVersion
	}

	var svcStartupFilePtr *valueObject.UnixFilePath
	if requestBody["startupFile"] != nil {
		svcStartupFile := valueObject.NewUnixFilePathPanic(
			requestBody["startupFile"].(string),
		)
		svcStartupFilePtr = &svcStartupFile
	}

	var svcPortBindings []valueObject.PortBinding
	if requestBody["portBindings"] != nil {
		svcPortBindings = parsePortBindings(
			requestBody["portBindings"].([]interface{}),
		)
	}

	autoCreateMapping := true
	if requestBody["autoCreateMapping"] != nil {
		var err error
		autoCreateMapping, err = sharedHelper.ParseBoolParam(
			requestBody["autoCreateMapping"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "InvalidAutoCreateMapping",
			)
		}
	}

	createInstallableServiceDto := dto.NewCreateInstallableService(
		svcName,
		svcVersionPtr,
		svcStartupFilePtr,
		svcPortBindings,
		autoCreateMapping,
	)

	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	servicesCmdRepo := servicesInfra.NewServicesCmdRepo()
	mappingQueryRepo := mappingInfra.NewMappingQueryRepo(controller.persistentDbSvc)
	mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)

	err := useCase.CreateInstallableService(
		servicesQueryRepo,
		servicesCmdRepo,
		mappingQueryRepo,
		mappingCmdRepo,
		vhostQueryRepo,
		createInstallableServiceDto,
	)

	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "InstallableServiceCreated")
}

// CreateCustomService godoc
// @Summary      CreateCustomService
// @Description  Install a new custom service.
// @Tags         services
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createCustomServiceDto	body dto.CreateCustomService	true	"name, type and command is required.<br />If version is not provided, it will be 'lts'.<br />If portBindings is not provided, it wil be default service port bindings.<br />If autoCreateMapping is not provided, it will be 'true'."
// @Success      201 {object} object{} "CustomServiceCreated"
// @Router       /v1/services/custom/ [post]
func (controller *ServicesController) CreateCustom(c echo.Context) error {
	requiredParams := []string{"name", "type", "command"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	svcName := valueObject.NewServiceNamePanic(requestBody["name"].(string))
	svcType := valueObject.NewServiceTypePanic(requestBody["type"].(string))
	svcCommand := valueObject.NewUnixCommandPanic(requestBody["command"].(string))

	var svcVersionPtr *valueObject.ServiceVersion
	if requestBody["version"] != nil {
		svcVersion := valueObject.NewServiceVersionPanic(
			requestBody["version"].(string),
		)
		svcVersionPtr = &svcVersion
	}

	var svcPortBindings []valueObject.PortBinding
	if requestBody["portBindings"] != nil {
		svcPortBindings = parsePortBindings(
			requestBody["portBindings"].([]interface{}),
		)
	}

	autoCreateMapping := true
	if requestBody["autoCreateMapping"] != nil {
		var err error
		autoCreateMapping, err = sharedHelper.ParseBoolParam(
			requestBody["autoCreateMapping"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "InvalidAutoCreateMapping",
			)
		}
	}

	createCustomServiceDto := dto.NewCreateCustomService(
		svcName,
		svcType,
		svcCommand,
		svcVersionPtr,
		svcPortBindings,
		autoCreateMapping,
	)

	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	servicesCmdRepo := servicesInfra.NewServicesCmdRepo()
	mappingQueryRepo := mappingInfra.NewMappingQueryRepo(controller.persistentDbSvc)
	mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)

	err := useCase.CreateCustomService(
		servicesQueryRepo,
		servicesCmdRepo,
		mappingQueryRepo,
		mappingCmdRepo,
		vhostQueryRepo,
		createCustomServiceDto,
	)

	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "CustomServiceCreated")
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
	requiredParams := []string{"name"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	svcName := valueObject.NewServiceNamePanic(requestBody["name"].(string))

	var svcTypePtr *valueObject.ServiceType
	if requestBody["type"] != nil {
		svcType := valueObject.NewServiceTypePanic(
			requestBody["type"].(string),
		)
		svcTypePtr = &svcType
	}

	var svcCommandPtr *valueObject.UnixCommand
	if requestBody["command"] != nil {
		svcCommand := valueObject.NewUnixCommandPanic(
			requestBody["command"].(string),
		)
		svcCommandPtr = &svcCommand
	}

	var svcStatusPtr *valueObject.ServiceStatus
	if requestBody["status"] != nil {
		svcStatus := valueObject.NewServiceStatusPanic(
			requestBody["status"].(string),
		)
		svcStatusPtr = &svcStatus
	}

	var svcVersionPtr *valueObject.ServiceVersion
	if requestBody["version"] != nil {
		svcVersion := valueObject.NewServiceVersionPanic(
			requestBody["version"].(string),
		)
		svcVersionPtr = &svcVersion
	}

	var svcStartupFilePtr *valueObject.UnixFilePath
	if requestBody["startupFile"] != nil {
		svcStartupFile := valueObject.NewUnixFilePathPanic(
			requestBody["startupFile"].(string),
		)
		svcStartupFilePtr = &svcStartupFile
	}

	var svcPortBindings []valueObject.PortBinding
	if requestBody["portBindings"] != nil {
		svcPortBindings = parsePortBindings(
			requestBody["portBindings"].([]interface{}),
		)
	}

	updateSvcDto := dto.NewUpdateService(
		svcName,
		svcTypePtr,
		svcCommandPtr,
		svcStatusPtr,
		svcVersionPtr,
		svcStartupFilePtr,
		svcPortBindings,
	)

	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	servicesCmdRepo := servicesInfra.NewServicesCmdRepo()
	mappingQueryRepo := mappingInfra.NewMappingQueryRepo(controller.persistentDbSvc)
	mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)

	err := useCase.UpdateService(
		servicesQueryRepo,
		servicesCmdRepo,
		mappingQueryRepo,
		mappingCmdRepo,
		updateSvcDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "ServiceUpdated")
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
	svcName := valueObject.NewServiceNamePanic(c.Param("svcName"))

	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	servicesCmdRepo := servicesInfra.NewServicesCmdRepo()
	mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)

	err := useCase.DeleteService(
		servicesQueryRepo,
		servicesCmdRepo,
		mappingCmdRepo,
		svcName,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "ServiceDeleted")
}
