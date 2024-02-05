package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
)

// GetServices	 godoc
// @Summary      GetServices
// @Description  List installed services and their status.
// @Tags         services
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Success      200 {array} dto.ServiceWithMetrics
// @Router       /services/ [get]
func GetServicesController(c echo.Context) error {
	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	servicesList, err := useCase.GetServicesWithMetrics(servicesQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, servicesList)
}

// GetServices	 godoc
// @Summary      GetInstallableServices
// @Description  List installable services.
// @Tags         services
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Success      200 {array} entity.InstallableService
// @Router       /services/installables/ [get]
func GetInstallableServicesController(c echo.Context) error {
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
// @Param        addInstallableServiceDto	body dto.CreateInstallableService	true	"CreateInstallableService"
// @Success      201 {object} object{} "InstallableServiceCreated"
// @Router       /services/installables/ [post]
func CreateInstallableServiceController(c echo.Context) error {
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
		autoCreateMapping, err = apiHelper.ParseBoolParam(
			requestBody["autoCreateMapping"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "InvalidAutoCreateMapping",
			)
		}
	}

	addInstallableServiceDto := dto.NewCreateInstallableService(
		svcName,
		svcVersionPtr,
		svcStartupFilePtr,
		svcPortBindings,
		autoCreateMapping,
	)

	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	servicesCmdRepo := servicesInfra.ServicesCmdRepo{}
	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

	err := useCase.CreateInstallableService(
		servicesQueryRepo,
		servicesCmdRepo,
		vhostQueryRepo,
		vhostCmdRepo,
		addInstallableServiceDto,
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
// @Param        addCustomServiceDto	body dto.CreateCustomService	true	"CreateCustomService"
// @Success      201 {object} object{} "CustomServiceCreated"
// @Router       /services/custom/ [post]
func CreateCustomServiceController(c echo.Context) error {
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
		autoCreateMapping, err = apiHelper.ParseBoolParam(
			requestBody["autoCreateMapping"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "InvalidAutoCreateMapping",
			)
		}
	}

	addCustomServiceDto := dto.NewCreateCustomService(
		svcName,
		svcType,
		svcCommand,
		svcVersionPtr,
		svcPortBindings,
		autoCreateMapping,
	)

	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	servicesCmdRepo := servicesInfra.ServicesCmdRepo{}
	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

	err := useCase.CreateCustomService(
		servicesQueryRepo,
		servicesCmdRepo,
		vhostQueryRepo,
		vhostCmdRepo,
		addCustomServiceDto,
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
// @Param        updateServiceDto	body dto.UpdateService	true	"UpdateServiceDetails"
// @Success      200 {object} object{} "ServiceUpdated"
// @Router       /services/ [put]
func UpdateServiceController(c echo.Context) error {
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
	servicesCmdRepo := servicesInfra.ServicesCmdRepo{}
	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

	err := useCase.UpdateService(
		servicesQueryRepo,
		servicesCmdRepo,
		vhostQueryRepo,
		vhostCmdRepo,
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
// @Param        svcName path string true "ServiceName"
// @Success      200 {object} object{} "ServiceDeleted"
// @Router       /services/{svcName}/ [delete]
func DeleteServiceController(c echo.Context) error {
	svcName := valueObject.NewServiceNamePanic(c.Param("svcName"))

	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	servicesCmdRepo := servicesInfra.ServicesCmdRepo{}

	err := useCase.DeleteService(
		servicesQueryRepo,
		servicesCmdRepo,
		svcName,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "ServiceDeleted")
}
