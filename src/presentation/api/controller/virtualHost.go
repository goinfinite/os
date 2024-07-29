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
	"github.com/speedianet/os/src/presentation/service"
)

type VirtualHostController struct {
	persistentDbSvc    *internalDbInfra.PersistentDatabaseService
	virtualHostService *service.VirtualHostService
}

func NewVirtualHostController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *VirtualHostController {
	return &VirtualHostController{
		persistentDbSvc:    persistentDbSvc,
		virtualHostService: service.NewVirtualHostService(persistentDbSvc),
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
	return apiHelper.ServiceResponseWrapper(c, controller.virtualHostService.Read())
}

// CreateVirtualHost    godoc
// @Summary      CreateVirtualHost
// @Description  Create a new vhost.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createVirtualHostDto 	  body    dto.CreateVirtualHost  true  "Only hostname is required.<br />type may be 'top-level', 'subdomain', 'wildcard' or 'alias'. If is not provided, it will be 'top-level'. If type is 'alias', parentHostname it will be required."
// @Success      201 {object} object{} "VirtualHostCreated"
// @Router       /v1/vhosts/ [post]
func (controller *VirtualHostController) Create(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.virtualHostService.Create(requestBody),
	)
}

// DeleteVirtualHost godoc
// @Summary      DeleteVirtualHost
// @Description  Delete a vhost.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        hostname path string true "Hostname to delete"
// @Success      200 {object} object{} "VirtualHostDeleted"
// @Router       /v1/vhosts/{hostname}/ [delete]
func (controller *VirtualHostController) Delete(c echo.Context) error {
	requestBody := map[string]interface{}{
		"hostname": c.Param("hostname"),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.virtualHostService.Delete(requestBody),
	)
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
	return apiHelper.ServiceResponseWrapper(c, controller.virtualHostService.ReadWithMappings())
}

// CreateVirtualHostMapping godoc
// @Summary      CreateVirtualHostMapping
// @Description  Create a new vhost mapping.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createMappingDto	body dto.CreateMapping	true	"hostname, path and targetType are required.<br />matchPattern may be 'begins-with', 'contains', 'equals' or 'ends-with'. If is not provided, it will be 'begins-with'.<br />targetType may be 'url', 'service', 'response-code', 'inline-html' or 'static-files'. If targetType is 'url', targetHttpResponseCode may be provided. If is not provided, targetHttpResponseCode will be '200'. If targetType is 'response-code', targetHttpResponseCode may be provided. If is not provided, targetValue will be required.<br />targetValue must have the same value as the targetType requires."
// @Success      201 {object} object{} "MappingCreated"
// @Router       /v1/vhosts/mapping/ [post]
func (controller *VirtualHostController) CreateMapping(c echo.Context) error {
	requiredParams := []string{"hostname", "path", "targetType"}
	requestBody, _ := apiHelper.ReadRequestBody(c)

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
	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(controller.persistentDbSvc)

	err := useCase.CreateMapping(
		mappingQueryRepo,
		mappingCmdRepo,
		vhostQueryRepo,
		servicesQueryRepo,
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
// @Param        mappingId path uint true "MappingId to delete."
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
