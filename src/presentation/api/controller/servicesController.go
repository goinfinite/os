package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/infra"
	apiHelper "github.com/speedianet/sam/src/presentation/api/helper"
)

// AuthLogin godoc
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
