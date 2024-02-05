package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
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
	vhostsQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhostsList, err := useCase.GetVirtualHosts(vhostsQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, vhostsList)
}

// CreateVirtualHost    godoc
// @Summary      AddNewVirtualHost
// @Description  Add a new vhost.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        addVirtualHostDto 	  body    dto.CreateVirtualHost  true  "NewVirtualHost (only hostname is required)."
// @Success      201 {object} object{} "VirtualHostCreated"
// @Router       /vhosts/ [post]
func CreateVirtualHostController(c echo.Context) error {
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

	addVirtualHostDto := dto.NewCreateVirtualHost(
		hostname,
		vhostType,
		parentHostnamePtr,
	)

	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

	err := useCase.CreateVirtualHost(
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
// @Param        hostname path string true "Hostname"
// @Success      200 {object} object{} "VirtualHostDeleted"
// @Router       /vhosts/{hostname}/ [delete]
func DeleteVirtualHostController(c echo.Context) error {
	hostname := valueObject.NewFqdnPanic(c.Param("hostname"))

	vhostsQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhostsCmdRepo := vhostInfra.VirtualHostCmdRepo{}

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
	vhostsQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhostsList, err := useCase.GetVirtualHostsWithMappings(vhostsQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, vhostsList)
}

// CreateMapping godoc
// @Summary      CreateMapping
// @Description  Create a new vhost mapping.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        addMappingDto	body dto.CreateMapping	true	"hostname, path and targetType are required. If targetType is 'url', targetUrl is required and so on.<br />targetType may be 'service', 'url' or 'response-code'.<br />matchPattern may be 'begins-with', 'contains', 'equals', 'ends-with' or empty."
// @Success      201 {object} object{} "MappingCreated"
// @Router       /vhosts/mapping/ [post]
func CreateVirtualHostMappingController(c echo.Context) error {
	requiredParams := []string{"hostname", "path", "targetType"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	hostname := valueObject.NewFqdnPanic(requestBody["hostname"].(string))
	path := valueObject.NewMappingPathPanic(requestBody["path"].(string))
	targetType := valueObject.NewMappingTargetTypePanic(
		requestBody["targetType"].(string),
	)

	matchPattern := valueObject.NewMappingMatchPatternPanic("begins-with")
	if requestBody["matchPattern"] != nil {
		matchPattern = valueObject.NewMappingMatchPatternPanic(
			requestBody["matchPattern"].(string),
		)
	}

	var targetServicePtr *valueObject.ServiceName
	if requestBody["targetService"] != nil {
		targetService := valueObject.NewServiceNamePanic(
			requestBody["targetService"].(string),
		)
		targetServicePtr = &targetService
	}

	var targetUrlPtr *valueObject.Url
	if requestBody["targetUrl"] != nil {
		targetUrl := valueObject.NewUrlPanic(requestBody["targetUrl"].(string))
		targetUrlPtr = &targetUrl
	}

	var targetHttpResponseCodePtr *valueObject.HttpResponseCode
	if requestBody["targetHttpResponseCode"] != nil {
		targetHttpResponseCode := valueObject.NewHttpResponseCodePanic(
			requestBody["targetHttpResponseCode"],
		)
		targetHttpResponseCodePtr = &targetHttpResponseCode
	}

	addMappingDto := dto.NewCreateMapping(
		hostname,
		path,
		matchPattern,
		targetType,
		targetServicePtr,
		targetUrlPtr,
		targetHttpResponseCodePtr,
	)

	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}
	svcsQueryRepo := servicesInfra.ServicesQueryRepo{}

	err := useCase.CreateMapping(
		vhostQueryRepo,
		vhostCmdRepo,
		svcsQueryRepo,
		addMappingDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "MappingCreated")
}

// DeleteVirtualHost godoc
// @Summary      DeleteMapping
// @Description  Delete a vhost mapping.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        hostname path string true "Hostname"
// @Param        mappingId path uint true "MappingId"
// @Success      200 {object} object{} "MappingDeleted"
// @Router       /vhosts/mapping/{hostname}/{mappingId}/ [delete]
func DeleteVirtualHostMappingController(c echo.Context) error {
	hostname := valueObject.NewFqdnPanic(c.Param("hostname"))
	mappingId := valueObject.NewMappingIdPanic(c.Param("mappingId"))

	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

	err := useCase.DeleteMapping(
		vhostQueryRepo,
		vhostCmdRepo,
		hostname,
		mappingId,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "MappingDeleted")
}
