package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/infra"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
)

// GetVirtualHosts	 godoc
// @Summary      GetVirtualHosts
// @Description  List virtual hosts.
// @Tags         vhosts
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Success      200 {array} entity.VirtualHost
// @Router       /vhosts/ [get]
func GetVirtualHostsController(c echo.Context) error {
	vhostsQueryRepo := infra.VirtualHostQueryRepo{}
	vhostsList, err := useCase.GetVirtualHosts(vhostsQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, vhostsList)
}

// GetVirtualHostsWithMappings	 godoc
// @Summary      GetVirtualHostsWithMappings
// @Description  List virtual hosts with mappings.
// @Tags         vhosts
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Success      200 {array} dto.VirtualHostWithMappings
// @Router       /vhosts/mapping/ [get]
func GetVirtualHostsWithMappingsController(c echo.Context) error {
	vhostsQueryRepo := infra.VirtualHostQueryRepo{}
	vhostsList, err := useCase.GetVirtualHostsWithMappings(vhostsQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, vhostsList)
}
