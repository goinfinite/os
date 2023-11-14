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
// @Description  List services and their status.
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

// UpdateService godoc
// @Summary      UpdateServiceStatus
// @Description  Start, stop, install or uninstall a service.
// @Tags         services
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateSvcStatusDto	body dto.UpdateSvcStatus	true	"UpdateServiceStatusDetails"
// @Success      200 {object} object{} "ServiceStatusUpdated"
// @Router       /services/ [put]
func UpdateServiceController(c echo.Context) error {
	requiredParams := []string{"name", "status"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	svcName := valueObject.NewServiceNamePanic(requestBody["name"].(string))
	svcStatus := valueObject.NewServiceStatusPanic(requestBody["status"].(string))
	var svcVersionPtr *valueObject.ServiceVersion
	if requestBody["version"] != nil {
		svcVersion := valueObject.NewServiceVersionPanic(
			requestBody["version"].(string),
		)
		svcVersionPtr = &svcVersion
	}

	updateSvcStatusDto := dto.NewUpdateSvcStatus(svcName, svcStatus, svcVersionPtr)

	servicesQueryRepo := infra.ServicesQueryRepo{}
	servicesCmdRepo := infra.ServicesCmdRepo{}

	err := useCase.UpdateServiceStatus(
		servicesQueryRepo,
		servicesCmdRepo,
		updateSvcStatusDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "ServiceStatusUpdated")
}
