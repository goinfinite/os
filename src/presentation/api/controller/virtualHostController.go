package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	"github.com/speedianet/os/src/infra"
	infraHelper "github.com/speedianet/os/src/infra/helper"
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

// DeleteVirtualHost godoc
// @Summary      DeleteVirtualHost
// @Description  Delete a vhost.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        hostname path string true "VirtualHostHostname"
// @Success      200 {object} object{} "VirtualHostDeleted"
// @Router       /vhosts/{hostname}/ [delete]
func DeleteVirtualHostController(c echo.Context) error {
	hostname := valueObject.NewFqdnPanic(c.Param("hostname"))

	vhostsQueryRepo := infra.VirtualHostQueryRepo{}
	vhostsCmdRepo := infra.VirtualHostCmdRepo{}

	primaryHostname, err := infraHelper.GetPrimaryHostname()
	if err != nil {
		panic("PrimaryHostnameNotFound")
	}

	err = useCase.DeleteVirtualHost(
		vhostsQueryRepo,
		vhostsCmdRepo,
		primaryHostname,
		hostname,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "VirtualHostDeleted")
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
