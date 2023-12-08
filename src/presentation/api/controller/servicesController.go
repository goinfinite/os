package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	"github.com/speedianet/os/src/infra"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
)

// GetServices	 godoc
// @Summary      GetServices
// @Description  List installed services and their status.
// @Tags         services
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Success      200 {array} entity.Service
// @Router       /services/ [get]
func GetServicesController(c echo.Context) error {
	servicesQueryRepo := infra.ServicesQueryRepo{}
	servicesList, err := useCase.GetServices(servicesQueryRepo)
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
	servicesQueryRepo := infra.ServicesQueryRepo{}
	servicesList, err := useCase.GetInstallableServices(servicesQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, servicesList)
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

	var svcPorts []valueObject.NetworkPort
	if requestBody["ports"] != nil {
		for _, port := range requestBody["ports"].([]interface{}) {
			svcPort := valueObject.NewNetworkPortPanic(port)
			svcPorts = append(svcPorts, svcPort)
		}
	}

	updateSvcDto := dto.NewUpdateService(
		svcName,
		svcTypePtr,
		svcCommandPtr,
		svcStatusPtr,
		svcVersionPtr,
		svcStartupFilePtr,
		svcPorts,
	)

	servicesQueryRepo := infra.ServicesQueryRepo{}
	servicesCmdRepo := infra.ServicesCmdRepo{}

	err := useCase.UpdateService(
		servicesQueryRepo,
		servicesCmdRepo,
		updateSvcDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "ServiceUpdated")
}
