package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	mappingInfra "github.com/speedianet/os/src/infra/vhost/mapping"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
)

type VirtualHostController struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewVirtualHostController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *VirtualHostController {
	return &VirtualHostController{
		persistentDbSvc: persistentDbSvc,
	}
}

// ReadVirtualHosts	 godoc
// @Summary      ReadVirtualHosts
// @Description  List virtual hosts.
// @Tags         vhosts
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Success      200 {array} entity.VirtualHost
// @Router       /v1/vhosts/ [get]
func (controller *VirtualHostController) Read(c echo.Context) error {
	vhostsQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)
	vhostsList, err := useCase.ReadVirtualHosts(vhostsQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, vhostsList)
}

// CreateVirtualHost    godoc
// @Summary      CreateVirtualHost
// @Description  Create a new vhost.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createVirtualHostDto 	  body    dto.CreateVirtualHost  true  "NewVirtualHost (only hostname is required)."
// @Success      201 {object} object{} "VirtualHostCreated"
// @Router       /v1/vhosts/ [post]
func (controller *VirtualHostController) Create(c echo.Context) error {
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

	createVirtualHostDto := dto.NewCreateVirtualHost(
		hostname,
		vhostType,
		parentHostnamePtr,
	)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)
	vhostCmdRepo := vhostInfra.NewVirtualHostCmdRepo(controller.persistentDbSvc)

	err := useCase.CreateVirtualHost(
		vhostQueryRepo,
		vhostCmdRepo,
		createVirtualHostDto,
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
// @Router       /v1/vhosts/{hostname}/ [delete]
func (controller *VirtualHostController) Delete(c echo.Context) error {
	hostname := valueObject.NewFqdnPanic(c.Param("hostname"))

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)
	vhostCmdRepo := vhostInfra.NewVirtualHostCmdRepo(controller.persistentDbSvc)

	primaryVhost, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		panic("PrimaryVirtualHostNotFound")
	}

	err = useCase.DeleteVirtualHost(
		vhostQueryRepo,
		vhostCmdRepo,
		primaryVhost,
		hostname,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "VirtualHostDeleted")
}

// ReadVirtualHostsWithMappings	 godoc
// @Summary      ReadVirtualHostsWithMappings
// @Description  List virtual hosts with mappings.
// @Tags         vhosts
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Success      200 {array} dto.VirtualHostWithMappings
// @Router       /v1/vhosts/mapping/ [get]
func (controller *VirtualHostController) ReadWithMappings(c echo.Context) error {
	mappingQueryRepo := mappingInfra.NewMappingQueryRepo(controller.persistentDbSvc)

	vhostsWithMappings, err := useCase.ReadVirtualHostsWithMappings(
		mappingQueryRepo,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, vhostsWithMappings)
}

// CreateVirtualHostMapping godoc
// @Summary      CreateVirtualHostMapping
// @Description  Create a new vhost mapping.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createMappingDto	body dto.CreateMapping	true	"hostname, path and targetType are required. If targetType is 'url', targetUrl is required and so on.<br />targetType may be 'service', 'url' or 'response-code'.<br />matchPattern may be 'begins-with', 'contains', 'equals', 'ends-with' or empty."
// @Success      201 {object} object{} "MappingCreated"
// @Router       /v1/vhosts/mapping/ [post]
func (controller *VirtualHostController) CreateMapping(c echo.Context) error {
	requiredParams := []string{"hostname", "path", "targetType"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	hostname := valueObject.NewFqdnPanic(requestBody["hostname"].(string))
	path := valueObject.NewMappingPathPanic(requestBody["path"].(string))

	matchPattern := valueObject.NewMappingMatchPatternPanic("begins-with")
	if requestBody["matchPattern"] != nil {
		matchPattern = valueObject.NewMappingMatchPatternPanic(
			requestBody["matchPattern"].(string),
		)
	}

	targetType := valueObject.NewMappingTargetTypePanic(
		requestBody["targetType"].(string),
	)

	var targetValuePtr *valueObject.MappingTargetValue
	if requestBody["targetValue"] != nil {
		targetValue := valueObject.NewMappingTargetValuePanic(
			requestBody["targetValue"], targetType,
		)
		targetValuePtr = &targetValue
	}

	var targetHttpResponseCodePtr *valueObject.HttpResponseCode
	if requestBody["targetHttpResponseCode"] != nil {
		targetHttpResponseCode := valueObject.NewHttpResponseCodePanic(
			requestBody["targetHttpResponseCode"],
		)
		targetHttpResponseCodePtr = &targetHttpResponseCode
	}

	createMappingDto := dto.NewCreateMapping(
		hostname,
		path,
		matchPattern,
		targetType,
		targetValuePtr,
		targetHttpResponseCodePtr,
	)

	mappingQueryRepo := mappingInfra.NewMappingQueryRepo(controller.persistentDbSvc)
	mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)
	svcsQueryRepo := servicesInfra.ServicesQueryRepo{}

	err := useCase.CreateMapping(
		mappingQueryRepo,
		mappingCmdRepo,
		vhostQueryRepo,
		svcsQueryRepo,
		createMappingDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "MappingCreated")
}

// DeleteVirtualHostMapping godoc
// @Summary      DeleteVirtualHostMapping
// @Description  Delete a vhost mapping.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        mappingId path uint true "MappingId"
// @Success      200 {object} object{} "MappingDeleted"
// @Router       /v1/vhosts/mapping/{mappingId}/ [delete]
func (controller *VirtualHostController) DeleteMapping(c echo.Context) error {
	mappingId := valueObject.NewMappingIdPanic(c.Param("mappingId"))

	mappingQueryRepo := mappingInfra.NewMappingQueryRepo(controller.persistentDbSvc)
	mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)

	err := useCase.DeleteMapping(
		mappingQueryRepo,
		mappingCmdRepo,
		mappingId,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "MappingDeleted")
}
