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

// AddVirtualHost    godoc
// @Summary      AddNewVirtualHost
// @Description  Add a new vhost.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        addVirtualHostDto 	  body    dto.AddVirtualHost  true  "NewVirtualHost"
// @Success      201 {object} object{} "VirtualHostCreated"
// @Router       /vhosts/ [post]
func AddVirtualHostController(c echo.Context) error {
	requiredParams := []string{"hostname"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	hostname := valueObject.NewFqdnPanic(requestBody["hostname"].(string))

	vhostTypeStr := "top-level"
	if requestBody["type"] != nil {
		vhostTypeStr = requestBody["type"].(string)
	}
	vhostType := valueObject.NewVirtualHostTypePanic(vhostTypeStr)

	var parentHostnamePtr *valueObject.Fqdn
	if requestBody["parentHostname"] != nil {
		parentHostname := valueObject.NewFqdnPanic(
			requestBody["parentHostname"].(string),
		)
		parentHostnamePtr = &parentHostname
	}

	addVirtualHostDto := dto.NewAddVirtualHost(
		hostname,
		vhostType,
		parentHostnamePtr,
	)

	vhostQueryRepo := infra.VirtualHostQueryRepo{}
	vhostCmdRepo := infra.VirtualHostCmdRepo{}

	err := useCase.AddVirtualHost(
		vhostQueryRepo,
		vhostCmdRepo,
		addVirtualHostDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "VirtualHostCreated")
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
